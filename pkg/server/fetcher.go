package server

import (
	"context"
	"fmt"
	"io"
	"multiplexer/pkg/log"
	"multiplexer/pkg/response"
	"net/http"
	"sync"
	"time"
)

type fetcher interface {
	Process(ctx context.Context, urls []string) ([]string, *response.ErrorResponse)
}

type Fetcher struct {
	client *http.Client
}

func NewFetcher(c *http.Client) *Fetcher {
	return &Fetcher{client: c}
}

func (f *Fetcher) Process(ctx context.Context, urls []string) ([]string, *response.ErrorResponse) {
	var wg sync.WaitGroup
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(requestTimeoutDuration))
	defer cancel()

	results := make(chan result, len(urls))
	sem := make(chan struct{}, maxOutRequests)

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			default:
				sem <- struct{}{}
				resp, err := f.fetchData(ctx, url)
				<-sem
				if err != nil {
					log.Infof("get error during resp: %s", err)
				} else {
					log.Infof("get resp")
				}
				results <- result{response: resp, err: err}
			}
		}(url)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var responseData []string
	for r := range results {
		if r.err != nil {
			log.Infof("cancel all requests due to error in response")
			cancel()
			return nil, &response.ErrorResponse{
				Error:  fmt.Sprintf("error processing URLs: %s", r.err.Error()),
				Status: http.StatusInternalServerError,
			}
		}
		responseData = append(responseData, string(r.response))
	}
	return responseData, nil
}

func (f *Fetcher) fetchData(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
