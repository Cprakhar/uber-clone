package main

import (
	"net/http"

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
func (s *httpServer) run() error {

	h := handler.NewHTTPHandler()

	srv := &http.Server{
		Addr:    s.addr,
		Handler: h,
	}

	return srv.ListenAndServe()
}
