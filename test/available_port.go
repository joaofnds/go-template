package test

import (
	"app/adapter/http"
	"net"
	"strconv"

	"go.uber.org/fx"
)

var (
	i                     int
	ports                 = FindPorts(10_000, 10)
	AvailablePortProvider = fx.Decorate(func(httpConfig http.Config) http.Config {
		httpConfig.Port = ports[i]
		i++
		return httpConfig
	})
)

func FindPorts(start, amount int) []int {
	ports := make([]int, 0, amount)

	for port := start; len(ports) < amount; port++ {
		if portIsAvailable(port) {
			ports = append(ports, port)
		}
	}

	return ports
}

func portIsAvailable(port int) bool {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return false
	}

	if err := l.Close(); err != nil {
		return false
	}

	return true
}
