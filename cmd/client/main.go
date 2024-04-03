package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"multiplexer/pkg/log"
	"multiplexer/pkg/response"
	"net/http"
	"time"
)

func post(url string, payload []byte, ch chan<- error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Millisecond*300))
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		cancel()
		return
	}
	req.Header.Set("Content-Type", "application/json")
	go func() {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Errorf("error during doing request: %s", err)
			ch <- err
			return
		}
		defer func() {
			if resp != nil {
				resp.Body.Close()
			}
		}()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			ch <- fmt.Errorf("error reading response body from %s: %v", url, err)
			return
		}

		log.Infof("Status code: %s\n", resp.Status)
		log.Infof("Response from %s: %s\n", url, body)

		ch <- nil // No error
	}()

	time.Sleep(time.Second)
	cancel()
}

func main() {
	urls := [][]string{
		{
			"http://example.com",
			"http://example.org",
			"http://example.net",
			"http://example.net",
			"http://example.net",
			"http://example.net",
			"http://example.net",
			"http://example.net",
			"http://example.net",
		},
	}

	errCh := make(chan error, len(urls))
	for _, p := range urls {
		payload, err := json.Marshal(response.MultiplexerRequest{Urls: p})
		if err != nil {
			panic(err)
		}
		post("http://localhost:8080/", payload, errCh)
	}

	for i := 0; i < len(urls); i++ {
		if err := <-errCh; err != nil {
			log.Infof("Error encountered: %v\n", err)
			break
		}
	}
}
