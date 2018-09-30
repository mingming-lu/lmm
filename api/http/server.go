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
	shutDownSignal := make(chan os.Signal, 1)
	signal.Notify(shutDownSignal, syscall.SIGINT, syscall.SIGTERM)
	return &Server{
		srv: &http.Server{
			Addr:    addr,
			Handler: router,
		},
		shutDownSignal: shutDownSignal,
		isShutDown:     false,
	}
}

// Run starts listening to this server and blocks goroutine
func (s *Server) Run() {
	if s.isShutDown {
		zap.L().Fatal("cannot run server since it is stopped")
	}

	zap.L().Info("start serving at " + s.srv.Addr)
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal(err.Error())
		}
	}()

	<-s.shutDownSignal
	close(s.shutDownSignal)

	s.ShutDown()
}

// ShutDown stops this server gracefully
func (s *Server) ShutDown() {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		zap.L().Fatal(err.Error())
	}

	s.isShutDown = true
	zap.L().Info("stopped serving at " + s.srv.Addr)
}
