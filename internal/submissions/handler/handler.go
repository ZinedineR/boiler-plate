package handler

import (
	"boiler-plate/internal/base/app"
	"boiler-plate/internal/base/handler"
	"boiler-plate/internal/submissions/domain"
	"boiler-plate/internal/submissions/service"
	"boiler-plate/pkg/exception"
	"boiler-plate/pkg/httputils"
	"boiler-plate/pkg/server"
	"net/http"
)

type HTTPHandler struct {
	App                *handler.BaseHTTPHandler
	SubmissionsService service.Service
}

func NewHTTPHandler(
	handler *handler.BaseHTTPHandler, SubmissionsService service.Service,
) *HTTPHandler {
	return &HTTPHandler{
		App:                handler,
		SubmissionsService: SubmissionsService,
	}
}

func (h HTTPHandler) Create(ctx *app.Context) *server.ResponseInterface {
	// Binding JSON
	request := domain.SubmissionRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		err := exception.InvalidArgument("error reading request")
		resException := httputils.GenErrorResponseException(err)
		return h.App.AsJsonInterface(ctx, http.StatusBadRequest, resException)
	}

	if err := h.SubmissionsService.Create(ctx, &request); err != nil {
		responseException := httputils.GenErrorResponseException(err)
		return h.App.AsJsonInterface(ctx, responseException.StatusCode, responseException)
	}
	return h.App.AsJsonInterface(ctx, http.StatusOK, httputils.DataSuccessResponse{
		StatusCode: http.StatusOK,
		Message:    "success created",
		Data:       request,
	})
}

func (h HTTPHandler) Detail(ctx *app.Context) *server.ResponseInterface {
	id := ctx.Param("id")

	// Exec Service
	detailAsset, errException := h.SubmissionsService.Detail(ctx, id)
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
	errException := h.SubmissionsService.Delete(ctx, id)
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
	result, err := h.SubmissionsService.Find(ctx, limitParam, pageParam)
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

func (h HTTPHandler) FindByUser(ctx *app.Context) *server.ResponseInterface {
	limitParam := ctx.DefaultQuery("pageSize", "0")
	pageParam := ctx.DefaultQuery("page", "0")
	id := ctx.Param("id")
	result, err := h.SubmissionsService.FindByUser(ctx, limitParam, pageParam, id)
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
