package req

import (
	"io"
	"net/http"
)

type Headers map[string]string

func Get(url string, headers Headers) (*http.Response, error) {
	return req(http.MethodGet, url, headers, nil)
}

func Post(url string, headers Headers, body io.Reader) (*http.Response, error) {
	return req(http.MethodPost, url, headers, body)
}

func Delete(url string, headers Headers) (*http.Response, error) {
	return req(http.MethodDelete, url, headers, nil)
}

func req(method string, url string, headers Headers, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return http.DefaultClient.Do(req)
}
