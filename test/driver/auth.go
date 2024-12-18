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
	url string
}

func NewAuthDriver(baseURL string) *AuthDriver {
	return &AuthDriver{baseURL}
}

func (d *AuthDriver) Login(username string, password string) (oauth2.Token, error) {
	var token oauth2.Token
	return token, makeJSONRequest(params{
		into:   &token,
		status: http.StatusOK,
		req: func() (*http.Response, error) {
			return req.Post(
				d.url+"/auth/login",
				map[string]string{"Content-Type": "application/json"},
				strings.NewReader(fmt.Sprintf(`{"username":%q,"password":%q}`, username, password)),
			)
		},
	})
}

func (d *AuthDriver) MustLogin(username string, password string) oauth2.Token {
	return matchers.Must2(d.Login(username, password))
}
