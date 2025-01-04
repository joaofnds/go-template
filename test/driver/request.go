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

func makeJSONRequest(params params) error {
	bytes, err := makeRequest(params.status, params.req)
	if err != nil {
		return err
	}

	if params.into == nil {
		return nil
	}

	return json.Unmarshal(bytes, params.into)
}

func makeRequest(status int, req func() (*http.Response, error)) ([]byte, error) {
	resp, err := req()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != status {
		return nil, RequestFailure{Status: resp.StatusCode, Body: string(b)}
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
