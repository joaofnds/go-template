package fiber

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	lblMethod = "method"
	lblPath   = "path"
	lblStatus = "status"
)

type PromInstrumentation struct {
	req *prometheus.CounterVec
}

type PromHabitInstrumentation struct{}

func NewPromHTTPInstrumentation() Instrumentation {
	return &PromInstrumentation{
		req: promauto.NewCounterVec(
			prometheus.CounterOpts{Name: "web_request"},
			[]string{lblMethod, lblPath, lblStatus},
		),
	}
}

func (i *PromInstrumentation) Middleware(ctx *fiber.Ctx) error {
	defer i.LogReq(ctx)
	return ctx.Next()
}

func (i *PromInstrumentation) LogReq(ctx *fiber.Ctx) {
	labels := prometheus.Labels{}
	labels[lblMethod] = string(ctx.Route().Method)
	labels[lblPath] = string(ctx.Route().Path)
	labels[lblStatus] = strconv.Itoa(ctx.Response().StatusCode())

	i.req.With(labels).Inc()
}
