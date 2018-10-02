package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

// Server is a wrapper of *http.Server and *Router
type Server struct {
	srv            *http.Server
	shutDownSignal chan os.Signal
	isShutDown     bool
}

// NewServer creates a new Server
func NewServer(addr string, router *Router) *Server {
	return &Server{
		srv: &http.Server{
			Addr:    addr,
			Handler: router,
		},
		shutDownSignal: make(chan os.Signal, 1),
		isShutDown:     false,
	}
}

// Run starts listening to this server and blocks goroutine
func (s *Server) Run() {
	if s.isShutDown {
		zap.L().Panic("cannot run server since it is stopped")
	}

	zap.L().Info("start serving at " + s.srv.Addr)
	go func() {
		signal.Notify(s.shutDownSignal, syscall.SIGINT, syscall.SIGTERM)

		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Panic(err.Error())
		}
	}()

	<-s.shutDownSignal
	s.ShutDown()
}

// ShutDown stops this server gracefully
func (s *Server) ShutDown() {
	signal.Stop(s.shutDownSignal)
	close(s.shutDownSignal)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		zap.L().Panic(err.Error())
	}

	s.isShutDown = true
	zap.L().Info("stopped serving at " + s.srv.Addr)
}
