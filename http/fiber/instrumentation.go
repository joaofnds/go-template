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

type PromHTTPInstrumentation struct {
	req *prometheus.CounterVec
}

type PromHabitInstrumentation struct{}

func NewPromHTTPInstrumentation() HTTPInstrumentation {
	return &PromHTTPInstrumentation{
		req: promauto.NewCounterVec(
			prometheus.CounterOpts{Name: "web_request"},
			[]string{lblMethod, lblPath, lblStatus},
		),
	}
}

func (i *PromHTTPInstrumentation) Middleware(ctx *fiber.Ctx) error {
	defer i.LogReq(ctx)
	return ctx.Next()
}

func (i *PromHTTPInstrumentation) LogReq(ctx *fiber.Ctx) {
	labels := prometheus.Labels{}
	labels[lblMethod] = string(ctx.Route().Method)
	labels[lblPath] = string(ctx.Route().Path)
	labels[lblStatus] = strconv.Itoa(ctx.Response().StatusCode())

	i.req.With(labels).Inc()
}
