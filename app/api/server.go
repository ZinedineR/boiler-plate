package api

import (
	"boiler-plate/app/appconf"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"boiler-plate/internal/base/handler"
	transHandler "boiler-plate/internal/transaction/handler"
	tempHandler "boiler-plate/internal/users/handler"
	"boiler-plate/pkg/server"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
)

type HttpServe struct {
	router                *gin.Engine
	base                  *handler.BaseHTTPHandler
	UsersHandler          *tempHandler.HTTPHandler
	TransactionHandler    *transHandler.HTTPHandler
	GRPCUserServer        *grpc.Server
	GRPCTransactionServer *grpc.Server
	GRPCGateway           *runtime.ServeMux
}

func (h *HttpServe) Run(config *appconf.Config) error {
	//h.setupUsersRouter()
	//h.setupDevRouter(config)
	h.setupGRPCUserRouter()
	h.setupGRPCTransactionRouter()
	h.base.Handlers = h
	//if h.base.IsStaging() {
	//	h.setupDevRouter()
	//}
	userConnection, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		if err := h.GRPCUserServer.Serve(userConnection); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	transactionConnection, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		if err := h.GRPCTransactionServer.Serve(transactionConnection); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	time.Sleep(time.Second)
	h.setupGRPCUserGateway()
	h.setupGRPCTransactionGateway()
	return h.router.Run(fmt.Sprintf(":%s", config.AppEnvConfig.HttpPort))
}

func New(
	appName string, base *handler.BaseHTTPHandler,
	Users *tempHandler.HTTPHandler, Trans *transHandler.HTTPHandler,
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
	grpcUserServer := grpc.NewServer()
	grpcTransactionServer := grpc.NewServer()
	grpcGatewayMux := runtime.NewServeMux()
	return &HttpServe{
		router:                r,
		base:                  base,
		UsersHandler:          Users,
		TransactionHandler:    Trans,
		GRPCUserServer:        grpcUserServer,
		GRPCTransactionServer: grpcTransactionServer,
		GRPCGateway:           grpcGatewayMux,
	}
}
