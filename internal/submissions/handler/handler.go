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

// AsJson always return httpStatus: 200, but Status field: 500,400,200...
func (h HTTPHandler) AsJson(ctx *app.Context, status int, message string, data interface{}) *server.Response {
	return h.App.AsJson(ctx, status, message, data)
}

func (h HTTPHandler) AsJsonInterface(ctx *app.Context, status int, data interface{}) *server.ResponseInterface {
	return h.App.AsJsonInterface(ctx, status, data)
}

// BadRequest For mobile, httpStatus:200, but Status field: http.MobileBadRequest
func (h HTTPHandler) BadRequest(ctx *app.Context, err error) *server.Response {
	return h.App.AsJson(ctx, http.StatusBadRequest, err.Error(), nil)
}

// ForbiddenRequest For mobile, httpStatus:200, but Status field: http.StatusForbidden
func (h HTTPHandler) ForbiddenRequest(ctx *app.Context, err error) *server.Response {
	return h.App.AsJson(ctx, http.StatusForbidden, err.Error(), nil)
}

func (h HTTPHandler) AsError(ctx *app.Context, message string) *server.Response {
	return h.App.AsJson(ctx, http.StatusInternalServerError, message, nil)
}

func (h HTTPHandler) ThrowBadRequestException(ctx *app.Context, message string) *server.Response {
	return h.App.ThrowExceptionJson(ctx, http.StatusBadRequest, 0, "Bad Request", message)
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
