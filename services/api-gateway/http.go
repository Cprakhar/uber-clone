package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cprakhar/uber-clone/services/api-gateway/handler"
)

type httpServer struct {
	addr string
}

// NewHTTPServer creates a new HTTP server instance
func NewHTTPServer(addr string) *httpServer {
	return &httpServer{addr: addr}
}

// run starts the HTTP server
func (s *httpServer) run() {

	h := handler.NewHTTPHandler()

	srv := &http.Server{
		Addr:    s.addr,
		Handler: h,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("server started listening on %s", s.addr)
		serverErrors <- srv.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Printf("error starting the server: %v", err)
	case sig := <-shutdown:
		log.Printf("server is shutting down due to signal: %v", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("could not gracefully shutdown the server: %v", err)
			srv.Close()
		}
	}
}
