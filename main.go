package main

import (
	"web/config"
	"web/http"
	"web/user"

	"go.uber.org/fx"
)

func main() {
	fx.New(config.Module, http.Module, user.Module).Run()
}
