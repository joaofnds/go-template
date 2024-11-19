package driver

import (
	"encoding/json"
	"io"
	"net/http"
)

type params struct {
	into   any
	status int
	req    func() (*http.Response, error)
}

func makeJSONRequest(p params) error {
	b, err := makeRequest(p.status, p.req)
	if err != nil {
		return err
	}

	if p.into == nil {
		return nil
	}

	return json.Unmarshal(b, p.into)
}

func makeRequest(status int, req func() (*http.Response, error)) ([]byte, error) {
	res, err := req()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != status {
		return nil, RequestFailure{Status: res.StatusCode, Body: string(b)}
	}

	return b, nil
}

type RequestFailure struct {
	Status int
	Body   string
}

func (e RequestFailure) Error() string {
	return e.Body
}
