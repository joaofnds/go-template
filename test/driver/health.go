package driver

import (
	"app/test/matchers"
	"app/test/req"

	"net/http"
)

type HealthDriver struct {
	url     string
	headers req.Headers
}

func NewHealthDriver(baseURL string, headers req.Headers) *HealthDriver {
	return &HealthDriver{url: baseURL + "/Health", headers: headers}
}

func (d *HealthDriver) GetReq() (*http.Response, error) {
	return req.Get(d.url, nil)
}

func (d *HealthDriver) MustGetReq() *http.Response {
	return matchers.Must2(d.GetReq())
}

func (d *HealthDriver) Get() (string, error) {
	b, err := makeRequest(
		http.StatusOK,
		func() (*http.Response, error) { return d.GetReq() },
	)

	return string(b), err
}

func (d *HealthDriver) MustGet() string {
	return matchers.Must2(d.Get())
}

func (d *HealthDriver) GetFailed() (string, error) {
	b, err := makeRequest(
		http.StatusServiceUnavailable,
		func() (*http.Response, error) { return d.GetReq() },
	)

	return string(b), err
}

func (d *HealthDriver) MustGetFailed() string {
	return matchers.Must2(d.GetFailed())
}
