package driver

import (
	"app/test/matchers"
	"app/test/req"
	"app/user"

	"fmt"
	"net/http"
	"strings"
)

type UserDriver struct {
	url     string
	headers req.Headers
}

func NewUserDriver(baseURL string, headers req.Headers) *UserDriver {
	return &UserDriver{url: baseURL}
}

func (driver *UserDriver) Create(name string) (user.User, error) {
	var u user.User

	return u, makeJSONRequest(params{
		into:   &u,
		status: http.StatusCreated,
		req: func() (*http.Response, error) {
			return req.Post(
				driver.url+"/users",
				req.MergeHeaders(driver.headers, map[string]string{"Content-Type": "application/json"}),
				strings.NewReader(fmt.Sprintf(`{"name":%q}`, name)),
			)
		},
	})
}

func (driver *UserDriver) MustCreate(name string) user.User {
	return matchers.Must2(driver.Create(name))
}

func (driver *UserDriver) Get(name string) (user.User, error) {
	var u user.User

	return u, makeJSONRequest(params{
		into:   &u,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				driver.url+"/users/"+name,
				req.MergeHeaders(driver.headers, map[string]string{"Content-Type": "application/json"}),
			)
		},
	})
}

func (driver *UserDriver) MustGet(name string) user.User {
	return matchers.Must2(driver.Get(name))
}

func (driver *UserDriver) List() ([]user.User, error) {
	var users []user.User
	return users, makeJSONRequest(params{
		into:   &users,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				driver.url+"/users",
				req.MergeHeaders(driver.headers, map[string]string{"Content-Type": "application/json"}),
			)
		},
	})
}

func (driver *UserDriver) MustList() []user.User {
	return matchers.Must2(driver.List())
}

func (driver *UserDriver) Delete(name string) error {
	return makeJSONRequest(params{
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Delete(
				driver.url+"/users/"+name,
				req.MergeHeaders(driver.headers, map[string]string{"Content-Type": "application/json"}),
			)
		},
	})
}

func (driver *UserDriver) MustDelete(name string) {
	matchers.Must(driver.Delete(name))
}

func (driver *UserDriver) GetFeature(name string) (map[string]any, error) {
	var features map[string]any

	return features, makeJSONRequest(params{
		into:   &features,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				driver.url+"/users/"+name+"/feature",
				req.MergeHeaders(driver.headers, map[string]string{"Content-Type": "application/json"}),
			)
		},
	})
}

func (driver *UserDriver) MustGetFeature(name string) map[string]any {
	return matchers.Must2(driver.GetFeature(name))
}
