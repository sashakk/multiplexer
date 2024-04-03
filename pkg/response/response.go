package response

type ErrorResponse struct {
	Error  string
	Status int
}

type MultiplexerRequest struct {
	Urls []string `json:"urls,omitempty"`
}

type MultiplexerResponse struct {
	Result []string `json:"result,omitempty"`
}
