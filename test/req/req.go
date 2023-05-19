package req

import (
	"io"
	"net/http"
)

var Default = New(http.DefaultClient)

func Get(url string, headers Headers) (*http.Response, error) {
	return Default.Get(url, headers)
}

func Post(url string, headers Headers, body io.Reader) (*http.Response, error) {
	return Default.Post(url, headers, body)
}

func Patch(url string, headers Headers, body io.Reader) (*http.Response, error) {
	return Default.Patch(url, headers, body)
}

func Delete(url string, headers Headers) (*http.Response, error) {
	return Default.Delete(url, headers)
}

type Headers map[string]string

type Req struct{ client *http.Client }

func New(client *http.Client) *Req {
	return &Req{client}
}

func (r *Req) Get(url string, headers Headers) (*http.Response, error) {
	return r.req(http.MethodGet, url, headers, nil)
}

func (r *Req) Post(url string, headers Headers, body io.Reader) (*http.Response, error) {
	return r.req(http.MethodPost, url, headers, body)
}

func (r *Req) Patch(url string, headers Headers, body io.Reader) (*http.Response, error) {
	return r.req(http.MethodPatch, url, headers, body)
}

func (r *Req) Delete(url string, headers Headers) (*http.Response, error) {
	return r.req(http.MethodDelete, url, headers, nil)
}

func (r *Req) req(method string, url string, headers Headers, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return r.client.Do(req)
}
