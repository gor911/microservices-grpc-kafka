package http

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func BuildHandler(h *Handlers) fasthttp.RequestHandler {
	r := router.New()

	r.GET("/health", h.Health)
	r.POST("/orders", h.CreateOrder)

	return r.Handler
}
