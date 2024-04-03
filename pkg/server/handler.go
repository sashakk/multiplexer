package server

import (
	"encoding/json"
	"fmt"
	"multiplexer/pkg/log"
	"multiplexer/pkg/response"
	"net/http"
)

type result struct {
	response []byte
	err      error
}

func (s *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {
	if !s.limiter.Allow() {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}
	defer s.limiter.Done()

	urls, errResponse := s.parseAndValidate(r)
	if errResponse != nil {
		http.Error(w, errResponse.Error, errResponse.Status)
		return
	}

	responses, errResponse := s.fetcher.Process(r.Context(), urls)
	if errResponse != nil {
		http.Error(w, errResponse.Error, errResponse.Status)
		return
	}

	jsonData, err := json.Marshal(response.MultiplexerResponse{Result: responses})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error processing URLs: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonData); err != nil {
		log.Errorf("Error while writing")
	}
}

func (s *Server) parseAndValidate(r *http.Request) ([]string, *response.ErrorResponse) {
	var request response.MultiplexerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, &response.ErrorResponse{Error: "Invalid JSON format", Status: http.StatusBadRequest}
	}

	if len(request.Urls) == 0 {
		return nil, &response.ErrorResponse{Error: "There is no urls to process", Status: http.StatusBadRequest}
	}

	if len(request.Urls) > 20 {
		return nil, &response.ErrorResponse{Error: "Too many URLs, maximum allowed is 20", Status: http.StatusBadRequest}
	}

	return request.Urls, nil
}
