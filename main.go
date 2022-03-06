package main

import (
	"web/config"
	"web/http"
	"web/logger"
	"web/user"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		logger.Module,
		config.Module,
		http.Module,
		user.Module,
	).Run()
}
