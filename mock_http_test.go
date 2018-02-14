package vertec

import (
	"net/http"
	"strings"
)

type MockRoundTripper struct {
	content string
	result	http.Response

}

type ReaderCloserType struct {
	reader *strings.Reader
}

func (this ReaderCloserType) Close() error {
	return nil
}

func (this *MockRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	this.result.Status = "200 OK"
	this.result.StatusCode = 200
	this.result.Body = ReaderCloserType{ strings.NewReader(this.content) }
	this.result.ContentLength = int64(len(this.content))
	return &this.result, nil
}

func (this ReaderCloserType) Read(p []byte) (n int, err error) {
	return this.reader.Read(p)
}
