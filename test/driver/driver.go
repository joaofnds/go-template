package driver

import (
	"app/user"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Driver struct {
	API *API
}

func NewDriver(url string) *Driver {
	return &Driver{API: NewAPI(url)}
}

func (d *Driver) CreateUser(name string) (user.User, error) {
	var u user.User

	return u, makeJSONRequest(params{
		into:   &u,
		status: http.StatusCreated,
		req:    func() (*http.Response, error) { return d.API.CreateUser(name) },
	})
}

func (d *Driver) GetUser(name string) (user.User, error) {
	var u user.User

	return u, makeJSONRequest(params{
		into:   &u,
		status: http.StatusOK,
		req:    func() (*http.Response, error) { return d.API.GetUser(name) },
	})
}

func (d *Driver) ListUsers() ([]user.User, error) {
	var users []user.User
	return users, makeJSONRequest(params{
		into:   &users,
		status: http.StatusOK,
		req:    d.API.ListUsers,
	})
}

func (d *Driver) DeleteUser(name string) error {
	return makeJSONRequest(params{
		status: http.StatusOK,
		req:    func() (*http.Response, error) { return d.API.DeleteUser(name) },
	})
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

	if res.StatusCode != p.status {
		return fmt.Errorf("expected status %d, got %d", p.status, res.StatusCode)
	}

	if p.into == nil {
		return nil
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, p.into)
}
