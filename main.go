package main

import (
	"web/config"
	"web/controllers"

	"go.uber.org/fx"
)

func main() {
	fx.New(config.Module, controllers.Module).Run()
}
