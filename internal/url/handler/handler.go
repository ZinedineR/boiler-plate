package handler

import (
	"boiler-plate/internal/base/handler"
	"boiler-plate/internal/url/domain"
	"boiler-plate/internal/url/service"
	"boiler-plate/pkg/exception"
	"boiler-plate/pkg/httputils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HTTPHandler struct {
	App         *handler.BaseHTTPHandler
	GRPCHandler *GRPCHandler
	URLService  service.Service
}

func NewHTTPHandler(
	handler *handler.BaseHTTPHandler, grpc *GRPCHandler, UsersService service.Service,
) *HTTPHandler {
	return &HTTPHandler{
		App:         handler,
		GRPCHandler: grpc,
		URLService:  UsersService,
	}
}

// Create handles the creation of a URL
// @Summary Create a new URL
// @Description Create a new URL and return the created URL
// @Tags URL
// @Accept json
// @Produce json
// @Param url body domain.URL true "URL"
// @Success 200 {object} httputils.DataSuccessResponse{data=domain.URL}
// @Failure 400 {object} httputils.ErrorResponse
// @Router /urls [post]
func (h HTTPHandler) Create(ctx *gin.Context) {
	// Binding JSON
	request := domain.URL{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		err := exception.InvalidArgument("error reading request")
		resException := httputils.GenErrorResponseException(err)
		h.App.AbortJSON(ctx, resException.StatusCode, resException)
	}

	if err := h.URLService.Create(ctx, &request); err != nil {
		responseException := httputils.GenErrorResponseException(err)
		h.App.AbortJSON(ctx, responseException.StatusCode, responseException)
	}
	h.App.JSON(ctx, http.StatusOK, httputils.DataSuccessResponse{
		StatusCode: http.StatusOK,
		Message:    "success created",
		Data:       request,
	})
}

func (h HTTPHandler) Detail(ctx *gin.Context) {
	id := ctx.Param("shorturl")

	// Exec Service
	detailAsset, errException := h.URLService.Detail(ctx, id)
	if errException != nil {
		respException := httputils.GenErrorResponseException(errException)
		h.App.AbortJSON(ctx, respException.StatusCode, respException)
	}
	h.App.Redirect(ctx, detailAsset.Longurl)
}
