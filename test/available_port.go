package test

import (
	"app/adapter/http"
	"app/adapter/metrics"
	"net"

	"go.uber.org/fx"
)

var (
	AvailablePortProvider = fx.Options(
		fx.Decorate(func(httpConfig http.Config) (http.Config, error) {
			addr, err := findAvailableAddr()
			if err != nil {
				return httpConfig, err
			}

			httpConfig.Port = addr.Port

			return httpConfig, nil
		}),
		fx.Decorate(func(metricsConfig metrics.Config) (metrics.Config, error) {
			addr, err := findAvailableAddr()
			if err != nil {
				return metricsConfig, err
			}

			metricsConfig.Addr = addr.String()

			return metricsConfig, nil
		}),
	)
)

func findAvailableAddr() (*net.TCPAddr, error) {
	listener, listenErr := net.Listen("tcp", ":0")
	if listenErr != nil {
		return nil, listenErr
	}

	if closeErr := listener.Close(); closeErr != nil {
		return nil, closeErr
	}

	return listener.Addr().(*net.TCPAddr), nil
}
