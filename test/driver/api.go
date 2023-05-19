package driver

import (
	"fmt"
	"net/http"
	"strings"

	"app/test/req"
)

type API struct {
	baseURL string
}

func NewAPI(baseURL string) *API {
	return &API{baseURL}
}

func (a API) CreateUser(name string) (*http.Response, error) {
	return req.Post(
		a.baseURL+"/users",
		map[string]string{"Content-Type": "application/json"},
		strings.NewReader(fmt.Sprintf(`{"name":%q}`, name)),
	)
}

func (a API) GetUser(name string) (*http.Response, error) {
	return req.Get(
		a.baseURL+"/users/"+name,
		map[string]string{"Accept": "application/json"},
	)
}

func (a API) ListUsers() (*http.Response, error) {
	return req.Get(
		a.baseURL+"/users",
		map[string]string{"Accept": "application/json"},
	)
}

func (a API) DeleteUser(name string) (*http.Response, error) {
	return req.Delete(
		a.baseURL+"/users/"+name,
		map[string]string{"Accept": "application/json"},
	)
}
