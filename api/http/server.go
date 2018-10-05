package http

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

var (
	isRestarted = flag.Bool("restart", false, "if is restarted process")
)

// Server is a wrapper of *http.Server and *Router
type Server struct {
	listener         net.Listener
	srv              *http.Server
	isShutDown       bool
	shutDownSignal   chan os.Signal
	stopReloadSignal chan struct{}
}

// NewServer creates a new Server
func NewServer(addr string, router *Router) *Server {
	flag.Parse()
	var (
		l   net.Listener
		err error
	)

	if *isRestarted {
		f := os.NewFile(3, fmt.Sprint(os.Getpid()))
		l, err = net.FileListener(f)
		defer f.Close()
	} else {
		l, err = net.Listen("tcp", addr)
	}

	if err != nil {
		log.Println("cannot listen addr: ", addr, err.Error())
	}

	return &Server{
		srv: &http.Server{
			Addr:    addr,
			Handler: router,
		},
		listener:   l,
		isShutDown: false,
	}
}

// Run starts listening to this server and blocks goroutine
func (s *Server) Run() {
	if s.isShutDown {
		zap.L().Panic("cannot run server since it is stopped")
	}

	s.shutDownSignal = make(chan os.Signal, 1)
	if *isRestarted {
		zap.L().Info("restart serving at "+s.srv.Addr,
			zap.Int("new pid", os.Getpid()),
			zap.Int("old pid", os.Getppid()),
		)
	} else {
		zap.L().Info("start serving at "+s.srv.Addr, zap.Int("pid", os.Getpid()))
	}
	go func() {
		signal.Notify(s.shutDownSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)

		if err := s.srv.Serve(s.listener); err != nil && err != http.ErrServerClosed {
			zap.L().Panic(err.Error())
		}
	}()

	s.stopReloadSignal = make(chan struct{}, 1)
	go func() {
		for {
			select {
			case <-s.stopReloadSignal:
				return
			default:
				lastHash := reloader.lastHash
				reloader.Recompile()
				reloader.CompareMD5AndSwap()
				if reloader.lastHash != lastHash {
					hangup()
				}
				<-time.After(5 * time.Second)
			}
		}
	}()

	s.ShutDown(<-s.shutDownSignal)
}

// ShutDown stops this server gracefully
func (s *Server) ShutDown(sig os.Signal) {
	s.stopReloadSignal <- struct{}{}
	signal.Stop(s.shutDownSignal)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if sig == syscall.SIGUSR2 {
		s.reload()
	}

	if err := s.srv.Shutdown(ctx); err != nil {
		zap.L().Panic(err.Error(), zap.Int("pid", os.Getpid()))
	}

	s.isShutDown = true
	zap.L().Info("stopped serving at "+s.srv.Addr, zap.Int("pid", os.Getpid()))
}

func hangup() {
	if err := syscall.Kill(os.Getpid(), syscall.SIGUSR2); err != nil {
		zap.L().Error("failed to kill process",
			zap.String("error", err.Error()),
			zap.Int("pid", os.Getpid()),
		)
	}
}

func (s *Server) reload() {
	l, ok := s.listener.(filer)
	if !ok {
		zap.L().Warn("listener is not a file")
		return
	}
	f, err := l.File()
	if err != nil {
		zap.L().Warn("cannot open listener's file descriptor")
	}
	defer f.Close()

	args := []string{"-restart"}
	cmd := exec.Command(os.Args[0], args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = []*os.File{f}
	cmd.Env = os.Environ()
	cmd.Run()
}

type filer interface {
	File() (*os.File, error)
}
