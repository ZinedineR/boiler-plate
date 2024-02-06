package handler

import (
	"github.com/google/uuid"
	"ms-batch/internal/base/app"
	BaseDomain "ms-batch/internal/base/domain"
	"ms-batch/internal/base/handler"
	"ms-batch/internal/template/domain"
	TemplateService "ms-batch/internal/template/service"
	"ms-batch/pkg/db"
	"ms-batch/pkg/responsehelper"
	"ms-batch/pkg/server"
	"net/http"
	"strconv"
	"strings"
)

type HTTPHandler struct {
	App             *handler.BaseHTTPHandler
	TemplateService TemplateService.Service
}

func NewHTTPHandler(
	handler *handler.BaseHTTPHandler, templateService TemplateService.Service,
) *HTTPHandler {
	return &HTTPHandler{
		App:             handler,
		TemplateService: templateService,
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

func (h HTTPHandler) CreateMailTemplate(ctx *app.Context) *server.ResponseInterface {
	body := domain.MailTemplate{ID: uuid.NewString()}
	request := domain.UpsertMailTemplate{ID: body.ID}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		respStatus := responsehelper.GetStatusResponse(http.StatusUnauthorized, "Error in reading request")
		return h.AsJsonInterface(ctx, http.StatusUnauthorized, respStatus)
	}

	duplicateCheck, err := h.TemplateService.FindOneMailTemplateByName(ctx, request.Name)
	if err != nil {
		respStatus := responsehelper.GetStatusResponse(http.StatusBadRequest, "Error in finding mail template")
		return h.AsJsonInterface(ctx, http.StatusBadRequest, respStatus)
	}

	if duplicateCheck != nil {
		respStatus := responsehelper.GetStatusResponse(http.StatusConflict, "found mail template with similar name in id: "+duplicateCheck.ID)
		return h.AsJsonInterface(ctx, http.StatusConflict, respStatus)
	}

	if err := h.TemplateService.CreateMailTemplate(ctx, &body); err != nil {
		respStatus := responsehelper.GetStatusResponse(http.StatusBadRequest, err.Error())
		return h.AsJsonInterface(ctx, http.StatusBadRequest, respStatus)
	}
	respStatus := responsehelper.GetStatusResponse(http.StatusOK, "")

	finalResponse := struct {
		*BaseDomain.Status
		*domain.MailTemplate
	}{respStatus, &body}
	return h.AsJsonInterface(ctx, http.StatusOK, finalResponse)
}

func (h HTTPHandler) FindMailTemplate(ctx *app.Context) *server.ResponseInterface {
	//body := domain.MailTemplate{ID: uuid.NewString()}
	var column, name string
	nameQuery := ctx.Query("name")
	if nameQuery != "" {
		filterQuery := strings.Split(nameQuery, ":")
		column = filterQuery[0]
		name = filterQuery[1]
	}
	limitQuery := ctx.Query("limit")
	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		if limitQuery == "" {
			limit = 0
		} else {
			respStatus := responsehelper.GetStatusResponse(http.StatusBadRequest, "Query must be in number")
			return h.AsJsonInterface(ctx, http.StatusBadRequest, respStatus)
		}
	}
	pageQuery := ctx.Query("page")
	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		if pageQuery == "" {
			page = 0
		} else {
			respStatus := responsehelper.GetStatusResponse(http.StatusBadRequest, "Query must be in number")
			return h.AsJsonInterface(ctx, http.StatusBadRequest, respStatus)
		}
	}
	result, pagination, err := h.TemplateService.FindMailTemplate(ctx, int64(limit), int64(page), column, name)
	if err != nil {
		respStatus := responsehelper.GetStatusResponse(http.StatusBadRequest, "Error in finding mail template")
		return h.AsJsonInterface(ctx, http.StatusBadRequest, respStatus)
	}

	respStatus := responsehelper.GetStatusResponse(http.StatusOK, "")

	finalResponse := struct {
		*BaseDomain.Status
		Pagination db.MongoPaginate      `json:"pagination"`
		Data       []domain.MailTemplate `json:"data"`
	}{respStatus, *pagination, *result}
	return h.AsJsonInterface(ctx, http.StatusOK, finalResponse)
}

func (h HTTPHandler) FindOneMailTemplate(ctx *app.Context) *server.ResponseInterface {
	idParam := ctx.Param("id")
	result, err := h.TemplateService.FindOneMailTemplate(ctx, idParam)
	if err != nil {
		respStatus := responsehelper.GetStatusResponse(http.StatusBadRequest, "Error in finding mail template")
		return h.AsJsonInterface(ctx, http.StatusBadRequest, respStatus)
	}

	if result == nil {
		respStatus := responsehelper.GetStatusResponse(http.StatusNotFound, "Can't find mail template with Id: "+idParam)
		return h.AsJsonInterface(ctx, http.StatusNotFound, respStatus)
	}

	respStatus := responsehelper.GetStatusResponse(http.StatusOK, "")

	finalResponse := struct {
		*BaseDomain.Status
		*domain.MailTemplate
	}{respStatus, result}
	return h.AsJsonInterface(ctx, http.StatusOK, finalResponse)
}

func (h HTTPHandler) UpdateMailTemplate(ctx *app.Context) *server.ResponseInterface {
	idParam := ctx.Param("id")
	body := domain.MailTemplate{ID: idParam}
	request := domain.UpsertMailTemplate{ID: idParam}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		respStatus := responsehelper.GetStatusResponse(http.StatusUnauthorized, "Error in reading json request")
		return h.AsJsonInterface(ctx, http.StatusUnauthorized, respStatus)
	}

	duplicateCheck, err := h.TemplateService.FindOneMailTemplateByName(ctx, request.Name)
	if err != nil {
		respStatus := responsehelper.GetStatusResponse(http.StatusBadRequest, "Error in finding mail template")
		return h.AsJsonInterface(ctx, http.StatusBadRequest, respStatus)
	}

	if duplicateCheck != nil {
		if duplicateCheck.ID != request.ID {
			respStatus := responsehelper.GetStatusResponse(http.StatusConflict, "found mail template with similar name in id: "+duplicateCheck.ID)
			return h.AsJsonInterface(ctx, http.StatusConflict, respStatus)
		}
	}
	if err := h.TemplateService.UpdateMailTemplate(ctx, idParam, &body); err != nil {
		respStatus := responsehelper.GetStatusResponse(http.StatusBadRequest, err.Error())
		return h.AsJsonInterface(ctx, http.StatusBadRequest, respStatus)
	}

	respStatus := responsehelper.GetStatusResponse(http.StatusOK, "")

	finalResponse := struct {
		*BaseDomain.Status
		*domain.MailTemplate
	}{respStatus, &body}
	return h.AsJsonInterface(ctx, http.StatusOK, finalResponse)
}

func (h HTTPHandler) DeleteMailTemplate(ctx *app.Context) *server.ResponseInterface {
	idParam := ctx.Param("id")
	if err := h.TemplateService.DeleteMailTemplate(ctx, idParam); err != nil {
		respStatus := responsehelper.GetStatusResponse(http.StatusBadRequest, "Error in creating mail template")
		return h.AsJsonInterface(ctx, http.StatusBadRequest, respStatus)
	}

	respStatus := responsehelper.GetStatusResponse(http.StatusOK, "")

	finalResponse := struct {
		*BaseDomain.Status
		AdditionalInfo string `json:"additional_info"`
	}{respStatus, idParam + " has been deleted"}
	return h.AsJsonInterface(ctx, http.StatusOK, finalResponse)
}
