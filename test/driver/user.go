package driver

import (
	"app/test/matchers"
	"app/test/req"
	"app/user"

	"net/http"
)

type UserDriver struct {
	url     string
	headers req.Headers
}

func NewUserDriver(baseURL string, headers req.Headers) *UserDriver {
	return &UserDriver{url: baseURL}
}

func (driver *UserDriver) Create(email string) (user.User, error) {
	var u user.User

	return u, makeJSONRequest(params{
		into:   &u,
		status: http.StatusCreated,
		req: func() (*http.Response, error) {
			return req.Post(
				driver.url+"/users",
				req.MergeHeaders(driver.headers, req.Headers{"Content-Type": "application/json"}),
				marshal(kv{"email": email}),
			)
		},
	})
}

func (driver *UserDriver) MustCreate(email string) user.User {
	return matchers.Must2(driver.Create(email))
}

func (driver *UserDriver) Get(userID string) (user.User, error) {
	var u user.User

	return u, makeJSONRequest(params{
		into:   &u,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				driver.url+"/users/"+userID,
				req.MergeHeaders(driver.headers, req.Headers{"Content-Type": "application/json"}),
			)
		},
	})
}

func (driver *UserDriver) MustGet(userID string) user.User {
	return matchers.Must2(driver.Get(userID))
}

func (driver *UserDriver) List() ([]user.User, error) {
	var users []user.User
	return users, makeJSONRequest(params{
		into:   &users,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				driver.url+"/users",
				req.MergeHeaders(driver.headers, req.Headers{"Content-Type": "application/json"}),
			)
		},
	})
}

func (driver *UserDriver) MustList() []user.User {
	return matchers.Must2(driver.List())
}

func (driver *UserDriver) Delete(userID string) error {
	return makeJSONRequest(params{
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Delete(
				driver.url+"/users/"+userID,
				req.MergeHeaders(driver.headers, req.Headers{"Content-Type": "application/json"}),
			)
		},
	})
}

func (driver *UserDriver) MustDelete(userID string) {
	matchers.Must(driver.Delete(userID))
}

func (driver *UserDriver) GetFeature(userID string) (map[string]any, error) {
	var features map[string]any

	return features, makeJSONRequest(params{
		into:   &features,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				driver.url+"/users/"+userID+"/feature",
				req.MergeHeaders(driver.headers, req.Headers{"Content-Type": "application/json"}),
			)
		},
	})
}

func (driver *UserDriver) MustGetFeature(userID string) map[string]any {
	return matchers.Must2(driver.GetFeature(userID))
}
