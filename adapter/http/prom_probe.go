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

type PromProbe struct {
	req *prometheus.CounterVec
}

type PromHabitProbe struct{}

func NewPromProbe() *PromProbe {
	return &PromProbe{
		req: promauto.NewCounterVec(
			prometheus.CounterOpts{Name: "app_request"},
			[]string{lblIP, lblMethod, lblPath, lblStatus},
		),
	}
}

func (p *PromProbe) Middleware(ctx *fiber.Ctx) error {
	defer p.LogReq(ctx)
	return ctx.Next()
}

func (p *PromProbe) LogReq(ctx *fiber.Ctx) {
	labels := prometheus.Labels{
		lblIP:     ctx.IP(),
		lblMethod: ctx.Route().Method,
		lblPath:   ctx.Route().Path,
		lblStatus: strconv.Itoa(ctx.Response().StatusCode()),
	}

	p.req.With(labels).Inc()
}
