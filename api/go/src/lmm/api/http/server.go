package http

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
		log.Panic("cannot run server since it is stopped")
	}

	log.Printf("start serving at %s, pid: %d\n", s.srv.Addr, os.Getpid())
	go func() {
		signal.Notify(s.shutDownSignal, syscall.SIGINT, syscall.SIGTERM)

		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panic(err)
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
		log.Panic(err)
	}

	s.isShutDown = true
	log.Printf("stopped serving at %s\n", s.srv.Addr)
}
