package handler

import (
	"boiler-plate/internal/base/app"
	BaseDomain "boiler-plate/internal/base/domain"
	"boiler-plate/internal/base/handler"
	"boiler-plate/internal/settings/domain"
	SettingsService "boiler-plate/internal/settings/service"
	"boiler-plate/pkg/responsehelper"
	"boiler-plate/pkg/server"
	"net/http"
)

type HTTPHandler struct {
	App             *handler.BaseHTTPHandler
	SettingsService SettingsService.Service
}

func NewHTTPHandler(
	handler *handler.BaseHTTPHandler, settingsService SettingsService.Service,
) *HTTPHandler {
	return &HTTPHandler{
		App:             handler,
		SettingsService: settingsService,
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

func (h HTTPHandler) FindSettings(ctx *app.Context) *server.ResponseInterface {
	result, err := h.SettingsService.FindSettings(ctx)
	if err != nil {
		respStatus := responsehelper.GetStatusResponse(http.StatusBadRequest, "Error in finding settings")
		return h.AsJsonInterface(ctx, http.StatusBadRequest, respStatus)
	}

	respStatus := responsehelper.GetStatusResponse(http.StatusOK, "")

	finalResponse := struct {
		*BaseDomain.Status
		*domain.MainTable
	}{respStatus, result}
	return h.AsJsonInterface(ctx, http.StatusOK, finalResponse)
}
