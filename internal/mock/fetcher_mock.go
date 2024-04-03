package mock

import (
	"context"
	"multiplexer/pkg/response"
)

type FetcherMock struct {
	ProcessFunc func(ctx context.Context, urls []string) ([]string, *response.ErrorResponse)
}

func (m *FetcherMock) Process(ctx context.Context, urls []string) ([]string, *response.ErrorResponse) {
	if m.ProcessFunc != nil {
		return m.ProcessFunc(ctx, urls)
	}
	return nil, nil
}
