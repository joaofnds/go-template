package authn

import "golang.org/x/oauth2"

type Provider interface {
	GetToken(username string, password string) (*oauth2.Token, error)
	ParseToken(token string) (Claims, error)
}

type Claims struct {
	Subject string `json:"sub,omitempty"`
}
