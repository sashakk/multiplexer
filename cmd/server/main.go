package main

import (
	"multiplexer/pkg/server"
	"net/http"
)

func main() {
	f := server.NewFetcher(&http.Client{})
	s := server.NewServer(f)
	s.Run()
}
