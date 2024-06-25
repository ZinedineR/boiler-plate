package api

import (
	"boiler-plate/app/appconf"
	"fmt"
	"os"
	"strings"

	"boiler-plate/internal/base/handler"
	subHandler "boiler-plate/internal/submissions/handler"
	tempHandler "boiler-plate/internal/users/handler"
	"boiler-plate/pkg/server"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
)

type HttpServe struct {
	router             *gin.Engine
	base               *handler.BaseHTTPHandler
	UsersHandler       *tempHandler.HTTPHandler
	SubmissionsHandler *subHandler.HTTPHandler
}

func (h *HttpServe) Run(config *appconf.Config) error {
	h.setupUsersRouter()
	h.setupDevRouter(config)
	h.base.Handlers = h

	//if h.base.IsStaging() {
	//	h.setupDevRouter()
	//}

	return h.router.Run(fmt.Sprintf(":%s", config.AppEnvConfig.HttpPort))
}

func New(
	appName string, base *handler.BaseHTTPHandler,
	Users *tempHandler.HTTPHandler,
	Submissions *subHandler.HTTPHandler,
) server.App {

	if os.Getenv("APP_ENV") != "production" {
		if os.Getenv("DEV_SHOW_ROUTE") == "False" {
			gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {

			}
		} else {
			gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
				fmt.Printf("Route: %-6s %-25s --> %s (%d handlers)\n",
					httpMethod, absolutePath, handlerName[strings.LastIndex(handlerName, "/")+1:], nuHandlers)

			}
		}
	}

	pathNamer := func(c *gin.Context) string {
		return fmt.Sprintf("%s %s%s", c.Request.Method, c.Request.Host, c.Request.RequestURI)
	}

	r := gin.New()
	r.Use(gintrace.Middleware(appName, gintrace.WithResourceNamer(pathNamer)))
	r.Use(ResponseHeaderFormat())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     base.AppConfig.AppEnvConfig.AllowOrigins,
		AllowMethods:     base.AppConfig.AppEnvConfig.AllowMethods,
		AllowHeaders:     base.AppConfig.AppEnvConfig.AllowHeaders,
		AllowCredentials: true,
	}))

	return &HttpServe{
		router:             r,
		base:               base,
		UsersHandler:       Users,
		SubmissionsHandler: Submissions,
	}
}
