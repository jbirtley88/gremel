package helper

import (
	"io"
	"net/http"
	"time"
)

// HttpHelper exists because there will be a plethora of different ways in which
// to interact with a HTTP endpoint, so we abstract it out to an interface
// and all of that complexity can be contained in one place to be refactored
// as patterns emerge.
//
// It also lends itself very well to mocking and stubbing for unit tests.
type HttpHelper interface {
	Get(url string) (code int, body io.ReadCloser, err error)
}

type HttpHelperBuilder struct {
	instance DefaultHttpHelper
}

func NewHttpHelperBuilder() *HttpHelperBuilder {
	return &HttpHelperBuilder{
		// Default timeout is 10 seconds
		instance: DefaultHttpHelper{
			timeout: time.Second * 10,
		},
	}
}

func (b *HttpHelperBuilder) WithTimeout(timeout time.Duration) *HttpHelperBuilder {
	b.instance.timeout = timeout
	return b
}

func (b *HttpHelperBuilder) Build() HttpHelper {
	defensiveCopy := b.instance
	return &defensiveCopy
}

type DefaultHttpHelper struct {
	timeout time.Duration
}

func (h *DefaultHttpHelper) Get(url string) (code int, body io.ReadCloser, err error) {
	client := http.Client{
		Timeout: h.timeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, resp.Body, nil
}
