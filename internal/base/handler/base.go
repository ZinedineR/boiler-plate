package handler

import (
	"boiler-plate/app/appconf"
	"boiler-plate/internal/base/app"
	baseModel "boiler-plate/pkg/db"
	"boiler-plate/pkg/httpclient"
	"boiler-plate/pkg/server"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HandlerFn func(ctx *app.Context) *server.Response
type HandlerFnInterface func(ctx *app.Context) *server.ResponseInterface

type BaseHTTPHandler struct {
	Handlers    interface{}
	DB          *gorm.DB
	AppConfig   *appconf.Config
	BaseModel   *baseModel.SQLClientRepository
	HttpClient  httpclient.Client
	GRPCHandler *GRPCHandler
}

func NewBaseHTTPHandler(
	db *gorm.DB,
	appConfig *appconf.Config,
	baseModel *baseModel.SQLClientRepository,
	httpClient httpclient.Client,
	grpcHandler *GRPCHandler,
) *BaseHTTPHandler {
	return &BaseHTTPHandler{
		DB:          db,
		AppConfig:   appConfig,
		BaseModel:   baseModel,
		HttpClient:  httpClient,
		GRPCHandler: grpcHandler,
	}
}

// Handler Basic Method ======================================================================================================

func (b BaseHTTPHandler) AsJsonInterface(ctx *app.Context, status int, data interface{}) *server.ResponseInterface {

	return &server.ResponseInterface{
		Status: status,
		Data:   data,
	}
}

// ThrowExceptionJson for some exception not handle in Yii2 framework
func (b BaseHTTPHandler) ThrowExceptionJson(ctx *app.Context, status, code int, name, message string) *server.Response {
	return &server.Response{
		Status:  status,
		Message: "",
		Log:     nil,
	}
}

func (b BaseHTTPHandler) UserAuthentication(c *gin.Context) (*app.Context, error) {
	return app.NewContext(c, b.AppConfig), nil
}

func (b BaseHTTPHandler) UserRunAction(handler HandlerFnInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		ctx, err := b.UserAuthentication(c)
		if err != nil {
			logrus.Errorln(fmt.Sprintf("REQUEST ID: %s , message: Unauthorized", ctx.APIReqID))
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Unauthorized",
				"data":    err.Error(),
			})
			return
		}

		defer func() {
			if err0 := recover(); err0 != nil {
				logrus.Errorln(err0)
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": "Request is halted unexpectedly, please contact the administrator.",
					"data":    nil,
				})
			}
		}()

		// Execute handler
		resp := handler(ctx)
		httpStatus := resp.Status

		if resp.Data == nil {
			c.Status(httpStatus)
			return
		}
		end := time.Now().Sub(start)
		logrus.Infoln(fmt.Sprintf("REQUEST ID: %s , LATENCY: %vms", ctx.APIReqID, end.Milliseconds()))
		c.JSON(httpStatus, resp.Data)

	}
}

func (b BaseHTTPHandler) GuestAuthentication(c *gin.Context) (*app.Context, error) {
	return app.NewContext(c, b.AppConfig), nil
}

func (b BaseHTTPHandler) GuestRunAction(handler HandlerFnInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		ctx, err := b.GuestAuthentication(c)
		if err != nil {
			logrus.Errorln(fmt.Sprintf("REQUEST ID: %s , message: Unauthorized", ctx.APIReqID))
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Unauthorized",
				"data":    err.Error(),
			})
			return
		}

		defer func() {
			if err0 := recover(); err0 != nil {
				logrus.Errorln(err0)
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": "Request is halted unexpectedly, please contact the administrator.",
					"data":    nil,
				})
			}
		}()

		resp := handler(ctx)
		httpStatus := resp.Status

		if resp.Data == nil {
			c.Status(httpStatus)
			return
		}
		end := time.Now().Sub(start)
		logrus.Infoln(fmt.Sprintf("REQUEST ID: %s , LATENCY: %vms", ctx.APIReqID, end.Milliseconds()))
		c.JSON(httpStatus, resp.Data)

	}
}
