package driver

import (
	"app/test/matchers"
	"app/test/req"

	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
)

type AuthDriver struct {
	url     string
	headers req.Headers
}

func NewAuthDriver(baseURL string, headers req.Headers) *AuthDriver {
	return &AuthDriver{url: baseURL, headers: headers}
}

func (driver *AuthDriver) Login(username string, password string) (oauth2.Token, error) {
	var token oauth2.Token
	return token, makeJSONRequest(params{
		into:   &token,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Post(
				driver.url+"/auth/login",
				map[string]string{"Content-Type": "application/json"},
				strings.NewReader(fmt.Sprintf(`{"username":%q,"password":%q}`, username, password)),
			)
		},
	})
}

func (driver *AuthDriver) MustLogin(username string, password string) oauth2.Token {
	token := matchers.Must2(driver.Login(username, password))
	driver.headers.Set("Authorization", fmt.Sprintf("%s %s", token.TokenType, token.AccessToken))
	return token
}

func (driver *AuthDriver) Register(username string, password string) error {
	_, err := makeRequest(
		http.StatusCreated,
		func() (*http.Response, error) {
			return req.Post(
				driver.url+"/auth/register",
				req.MergeHeaders(driver.headers, req.Headers{"Content-Type": "application/json"}),
				strings.NewReader(fmt.Sprintf(`{"username":%q,"password":%q}`, username, password)),
			)
		},
	)

	return err
}

func (driver *AuthDriver) MustRegister(username string, password string) {
	matchers.Must(driver.Register(username, password))
}

func (driver *AuthDriver) UserInfo() (map[string]string, error) {
	var userInfo map[string]string
	return userInfo, makeJSONRequest(params{
		into:   &userInfo,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Get(
				driver.url+"/auth/userinfo",
				req.MergeHeaders(driver.headers, req.Headers{"Content-Type": "application/json"}),
			)
		},
	})
}

func (driver *AuthDriver) MustUserInfo() map[string]string {
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
