package handler

import (
	"boiler-plate/internal/base/handler"
	"boiler-plate/internal/transaction/service"
)

type HTTPHandler struct {
	App                *handler.BaseHTTPHandler
	GRPCHandler        *GRPCHandler
	TransactionService service.Service
}

func NewHTTPHandler(
	handler *handler.BaseHTTPHandler, grpc *GRPCHandler, TransactionService service.Service,
) *HTTPHandler {
	return &HTTPHandler{
		App:                handler,
		GRPCHandler:        grpc,
		TransactionService: TransactionService,
	}
}
