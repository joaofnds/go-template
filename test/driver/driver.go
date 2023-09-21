package driver

import (
	apphttp "app/adapter/http"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/fx"
)

var Provider = fx.Provide(func(config apphttp.Config) *Driver {
	return NewDriver(fmt.Sprintf("http://localhost:%d", config.Port))
})

type Driver struct {
	URL  string
	User *UserDriver
}

func NewDriver(url string) *Driver {
	return &Driver{
		URL:  url,
		User: NewUserDriver(url),
	}
}

type params struct {
	into   any
	status int
	req    func() (*http.Response, error)
}

func makeJSONRequest(p params) error {
	res, err := p.req()
	if err != nil {
		return err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != p.status {
		return RequestFailure{Status: res.StatusCode, Body: string(b)}
	}

	if p.into == nil {
		return nil
	}

	return json.Unmarshal(b, p.into)
}

type RequestFailure struct {
	Status int
	Body   string
}

func (e RequestFailure) Error() string {
	return e.Body
}
