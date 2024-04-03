package server

import (
	"context"
	"fmt"
	"multiplexer/pkg/limiter"
	"multiplexer/pkg/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	maxOutRequests         = 4
	maxInRequests          = 100
	requestTimeoutDuration = time.Second
	host                   = ""
	port                   = "8080"
)

type Server struct {
	limiter *limiter.RequestLimiter
	fetcher fetcher
}

func NewServer(f fetcher) *Server {
	return &Server{
		limiter: limiter.NewRequestLimiter(maxInRequests),
		fetcher: f,
	}
}

func (s *Server) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: mux,
	}
	mux.Handle("/", withRecovery(validationMiddleware(http.HandlerFunc(s.HandleRequest))))

	go func() {
		log.Infof("start serving...")
		if err := server.ListenAndServe(); err != nil {
			log.Errorf("error during serve: %s", err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals

	log.Infof("Shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("Error shutting down server: %v\n", err)
	}

	log.Infof("Server shutdown complete")
}
