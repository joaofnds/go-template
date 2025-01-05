package driver

import (
	"app/test/matchers"
	"app/test/req"
	"app/user"

	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

type AuthDriver struct {
	url     string
	headers req.Headers
}

func NewAuthDriver(baseURL string, headers req.Headers) *AuthDriver {
	return &AuthDriver{url: baseURL, headers: headers}
}

func (driver *AuthDriver) Login(email string, password string) (oauth2.Token, error) {
	var token oauth2.Token
	return token, makeJSONRequest(params{
		into:   &token,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Post(
				driver.url+"/auth/login",
				req.Headers{"Content-Type": "application/json"},
				marshal(kv{"email": email, "password": password}),
			)
		},
	})
}

func (driver *AuthDriver) MustLogin(email string, password string) oauth2.Token {
	return matchers.Must2(driver.Login(email, password))
}

func (driver *AuthDriver) Register(email string, password string) (user.User, error) {
	var user user.User
	return user, makeJSONRequest(params{
		into:   &user,
		status: http.StatusCreated,
		req: func() (*http.Response, error) {
			return req.Post(
				driver.url+"/auth/register",
				req.MergeHeaders(driver.headers, req.Headers{"Content-Type": "application/json"}),
				marshal(kv{"email": email, "password": password}),
			)
		},
	})
}

func (driver *AuthDriver) MustRegister(email string, password string) user.User {
	return matchers.Must2(driver.Register(email, password))
}

func (driver *AuthDriver) UserInfo() (user.User, error) {
	var u user.User
	return u, makeJSONRequest(params{
		into:   &u,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(driver.url+"/auth/userinfo", driver.headers)
		},
	})
}

func (driver *AuthDriver) MustUserInfo() user.User {
	return matchers.Must2(driver.UserInfo())
}

func (driver *AuthDriver) Delete(email string) error {
	reqURL := matchers.Must2(url.Parse(driver.url))
	reqURL.Path = "/auth/delete"
	reqURL.RawQuery = url.Values{"email": {email}}.Encode()

	_, err := makeRequest(
		http.StatusNoContent,
		func() (*http.Response, error) {
			return req.Delete(reqURL.String(), driver.headers)
		},
	)

	return err
}

func (driver *AuthDriver) MustDelete(email string) {
	matchers.Must(driver.Delete(email))
}
