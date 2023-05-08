package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	lblIP     = "ip"
	lblMethod = "method"
	lblPath   = "path"
	lblStatus = "status"
)

type PromInstrumentation struct {
	req *prometheus.CounterVec
}

type PromHabitInstrumentation struct{}

func NewPromHTTPInstrumentation() *PromInstrumentation {
	return &PromInstrumentation{
		req: promauto.NewCounterVec(
			prometheus.CounterOpts{Name: "app_request"},
			[]string{lblIP, lblMethod, lblPath, lblStatus},
		),
	}
}

func (i *PromInstrumentation) Middleware(ctx *fiber.Ctx) error {
	defer i.LogReq(ctx)
	return ctx.Next()
}

func (i *PromInstrumentation) LogReq(ctx *fiber.Ctx) {
	labels := prometheus.Labels{
		lblIP:     ctx.IP(),
		lblMethod: ctx.Route().Method,
		lblPath:   ctx.Route().Path,
		lblStatus: strconv.Itoa(ctx.Response().StatusCode()),
	}

	i.req.With(labels).Inc()
}
