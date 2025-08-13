package http

import (
	"encoding/json"

	"github.com/gor911/microservices-grpc-kafka/order-service/internal/dto"
	"github.com/gor911/microservices-grpc-kafka/order-service/internal/service"
	"github.com/valyala/fasthttp"
)

type Handlers struct {
	orderService *service.Order
}

func NewHandlers(orderService *service.Order) *Handlers {
	return &Handlers{orderService}
}

func (h *Handlers) Health(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (h *Handlers) CreateOrder(ctx *fasthttp.RequestCtx) {
	input := dto.CreateOrderInput{}

	err := json.Unmarshal(ctx.PostBody(), &input)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)

		return
	}

	output, err := h.orderService.CreateOrder(ctx, input)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusUnprocessableEntity)

		return
	}

	jsonBytes, err := json.Marshal(output)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)

		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(jsonBytes)
}
