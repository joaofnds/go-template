package driver

import (
	"app/test/req"
	"app/user"

	"fmt"
	"net/http"
	"strings"
)

type UserDriver struct {
	url string
}

func NewUserDriver(baseURL string) *UserDriver {
	return &UserDriver{baseURL}
}

func (d *UserDriver) CreateUser(name string) (user.User, error) {
	var u user.User

	return u, makeJSONRequest(params{
		into:   &u,
		status: http.StatusCreated,
		req: func() (*http.Response, error) {
			return req.Post(
				d.url+"/users",
				map[string]string{"Content-Type": "application/json"},
				strings.NewReader(fmt.Sprintf(`{"name":%q}`, name)),
			)
		},
	})
}

func (d *UserDriver) GetUser(name string) (user.User, error) {
	var u user.User

	return u, makeJSONRequest(params{
		into:   &u,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				d.url+"/users/"+name,
				map[string]string{"Accept": "application/json"},
			)
		},
	})
}

func (d *UserDriver) ListUsers() ([]user.User, error) {
	var users []user.User
	return users, makeJSONRequest(params{
		into:   &users,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				d.url+"/users",
				map[string]string{"Accept": "application/json"},
			)
		},
	})
}

func (d *UserDriver) DeleteUser(name string) error {
	return makeJSONRequest(params{
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Delete(
				d.url+"/users/"+name,
				map[string]string{"Accept": "application/json"},
			)
		},
	})
}

func (d *UserDriver) GetFeature(name string) (map[string]any, error) {
	var features map[string]any

	return features, makeJSONRequest(params{
		into:   &features,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				d.url+"/users/"+name+"/feature",
				map[string]string{"Accept": "application/json"},
			)
		},
	})
}
