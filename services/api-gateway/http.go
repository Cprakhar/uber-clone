package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cprakhar/uber-clone/services/api-gateway/handler"
	"github.com/cprakhar/uber-clone/shared/messaging"
	"github.com/cprakhar/uber-clone/shared/messaging/kafka"
)

type httpServer struct {
	addr        string
	kfClient    *kafka.KafkaClient
	connManager *messaging.ConnectionManager
}

// NewhttpServer creates a new http server instance
func NewhttpServer(addr string, kfc *kafka.KafkaClient, connMgr *messaging.ConnectionManager) *httpServer {
	return &httpServer{addr: addr, kfClient: kfc, connManager: connMgr}
}

// run starts the http server
func (s *httpServer) run(ctx context.Context) error {
	// http server setup
	h := handler.NewHTTPHandler(s.kfClient, s.connManager)
	srv := &http.Server{
		Addr:    s.addr,
		Handler: h,
	}

	// Start the server in a separate goroutine
	errCh := make(chan error, 1)
	go func() {
		log.Printf("http server running on %s", s.addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	// Wait for context cancellation or server error
	select {
	case <-ctx.Done():
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("http server error: %w", err)
		}
	}

	// Graceful shutdown with a timeout
	shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shCtx); err != nil {
		return fmt.Errorf("http server shutdown error: %w", err)
	}

	log.Println("http server gracefully stopped")
	return nil
}
