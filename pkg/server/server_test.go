package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"multiplexer/internal/mock"
	"multiplexer/pkg/response"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleRequest_ValidInput(t *testing.T) {
	mockFetcher := &mock.FetcherMock{
		ProcessFunc: func(ctx context.Context, urls []string) ([]string, *response.ErrorResponse) {
			return []string{}, nil
		},
	}

	server := NewServer(mockFetcher)

	request := response.MultiplexerRequest{}
	request.Urls = append(request.Urls, "http://localhost:8081/")
	data, err := json.Marshal(request)
	if err != nil {
		t.Errorf("error while marshal: %s", err)
	}

	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	server.HandleRequest(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestHandleRequest_MoreThan20URLs(t *testing.T) {
	mockFetcher := &mock.FetcherMock{
		ProcessFunc: func(ctx context.Context, urls []string) ([]string, *response.ErrorResponse) {
			return []string{}, nil
		},
	}

	server := NewServer(mockFetcher)

	request := response.MultiplexerRequest{}
	for i := 0; i < 21; i++ {
		request.Urls = append(request.Urls, "http://localhost:8081/")
	}
	data, err := json.Marshal(request)
	if err != nil {
		t.Errorf("error while marshal: %s", err)
	}

	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	server.HandleRequest(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandleRequest_Valid20URLs(t *testing.T) {
	shutdown := make(chan struct{})
	go runTestServer(shutdown)

	f := NewFetcher(&http.Client{})
	server := NewServer(f)
	request := response.MultiplexerRequest{}
	for i := 0; i < 20; i++ {
		request.Urls = append(request.Urls, "http://localhost:8081/")
	}
	data, err := json.Marshal(request)
	if err != nil {
		t.Errorf("error while marshal: %s", err)
	}
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	server.HandleRequest(rr, req)
	shutdown <- struct{}{}
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}
}

func runTestServer(shutdown chan struct{}) {
	mux := http.NewServeMux()
	testServer := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	go testServer.ListenAndServe()
	<-shutdown
	testServer.Shutdown(context.Background())
}
