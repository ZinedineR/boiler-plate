package handler

import (
	"boiler-plate/internal/base/app"
	"boiler-plate/internal/base/handler"
	"boiler-plate/internal/users/domain"
	"boiler-plate/internal/users/service"
	"boiler-plate/pkg/exception"
	"boiler-plate/pkg/httputils"
	"boiler-plate/pkg/server"
	"net/http"
)

type HTTPHandler struct {
	App          *handler.BaseHTTPHandler
	GRPCHandler  *GRPCHandler
	UsersService service.Service
}

func NewHTTPHandler(
	handler *handler.BaseHTTPHandler, grpc *GRPCHandler, UsersService service.Service,
) *HTTPHandler {
	return &HTTPHandler{
		App:          handler,
		GRPCHandler:  grpc,
		UsersService: UsersService,
	}
}

func (h HTTPHandler) Create(ctx *app.Context) *server.ResponseInterface {
	// Binding JSON
	request := domain.Users{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		err := exception.InvalidArgument("error reading request")
		resException := httputils.GenErrorResponseException(err)
		return h.App.AsJsonInterface(ctx, http.StatusBadRequest, resException)
	}

	//if err := h.UsersService.Create(ctx, &request); err != nil {
	//	responseException := httputils.GenErrorResponseException(err)
	//	return h.App.AsJsonInterface(ctx, responseException.StatusCode, responseException)
	//}
	return h.App.AsJsonInterface(ctx, http.StatusOK, httputils.DataSuccessResponse{
		StatusCode: http.StatusOK,
		Message:    "success created",
		Data:       request,
	})
}

func (h HTTPHandler) Update(ctx *app.Context) *server.ResponseInterface {
	//id := ctx.Param("id")
	// Binding JSON
	request := domain.Users{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		err := exception.InvalidArgument("error reading request")
		resException := httputils.GenErrorResponseException(err)
		return h.App.AsJsonInterface(ctx, http.StatusBadRequest, resException)
	}

	// Exec Service
	//errException := h.UsersService.Update(ctx, id, &request)
	//if errException != nil {
	//	responseException := httputils.GenErrorResponseException(errException)
	//	return h.App.AsJsonInterface(ctx, responseException.StatusCode, responseException)
	//}

	// return
	return h.App.AsJsonInterface(ctx, http.StatusOK, httputils.DataSuccessResponse{
		StatusCode: http.StatusOK,
		Message:    "success update",
		Data:       request,
	})
}

func (h HTTPHandler) Detail(ctx *app.Context) *server.ResponseInterface {
	id := ctx.Param("id")

	// Exec Service
	detailAsset, errException := h.UsersService.Detail(ctx, id)
	if errException != nil {
		respException := httputils.GenErrorResponseException(errException)
		return h.App.AsJsonInterface(ctx, respException.StatusCode, respException)
	}
	return h.App.AsJsonInterface(ctx, http.StatusOK, httputils.DataSuccessResponse{
		StatusCode: http.StatusOK,
		Message:    "success",
		Data:       detailAsset,
	})
}

func (h HTTPHandler) Delete(ctx *app.Context) *server.ResponseInterface {
	id := ctx.Param("id")

	// Exec Service
	errException := h.UsersService.Delete(ctx, id)
	if errException != nil {
		respException := httputils.GenErrorResponseException(errException)
		return h.App.AsJsonInterface(ctx, respException.StatusCode, respException)
	}
	return h.App.AsJsonInterface(ctx, http.StatusOK, httputils.SuccessResponse{
		StatusCode: http.StatusOK,
		Message:    "success delete id: " + id,
	})
}

func (h HTTPHandler) Find(ctx *app.Context) *server.ResponseInterface {
	limitParam := ctx.DefaultQuery("pageSize", "0")
	pageParam := ctx.DefaultQuery("page", "0")
	result, err := h.UsersService.Find(ctx, limitParam, pageParam)
	if err != nil {
		responseException := httputils.GenErrorResponseException(err)
		return h.App.AsJsonInterface(ctx, responseException.StatusCode, responseException)
	}

	return h.App.AsJsonInterface(ctx, http.StatusOK, httputils.DataSuccessResponse{
		StatusCode: http.StatusOK,
		Message:    "success",
		Data:       result,
	})
}
