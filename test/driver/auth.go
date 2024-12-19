package driver

import (
	"app/test/matchers"
	"app/test/req"

	"fmt"
	"net/http"
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
