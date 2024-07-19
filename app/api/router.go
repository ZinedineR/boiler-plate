package api

import (
	"boiler-plate/internal/base/handler"
	pbUrl "boiler-plate/proto/url/v1"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func (h *HttpServe) setupGRPCRouter() {
	pbUrl.RegisterServiceServer(h.GRPCServer, h.URLHandler.GRPCHandler)

}

func (h *HttpServe) setupURLRouter() {
	urlGroup := h.router.Group("api/v2/url")
	urlGroup.POST("/", h.URLHandler.Create)
	urlGroup.GET("/:shorturl", h.URLHandler.Detail)
}

func (h *HttpServe) setupGRPCGateway() {
	client, err := grpc.NewClient(
		"0.0.0.0:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	if err := pbUrl.RegisterServiceHandler(context.Background(), h.GRPCGateway, client); err != nil {
		log.Fatal("failed to register users gateway")
	}
	h.router.Group("/api/v2/*{grpc_gateway}").Any("", gin.WrapH(h.GRPCGateway))
}
func (h *HttpServe) UserRoute(method, path string, f handler.HandlerFnInterface) {
	userRoute := h.router.Group("/api/v2")
	switch method {
	case "GET":
		userRoute.GET(path, h.base.UserRunAction(f))
	case "POST":
		userRoute.POST(path, h.base.UserRunAction(f))
	case "PUT":
		userRoute.PUT(path, h.base.UserRunAction(f))
	case "DELETE":
		userRoute.DELETE(path, h.base.UserRunAction(f))
	default:
		panic(fmt.Sprintf(":%s method not allow", method))
	}
}

func (h *HttpServe) GuestRoute(method, path string, f handler.HandlerFnInterface) {
	guestRoute := h.router.Group("/api/v2")
	switch method {
	case "GET":
		guestRoute.GET(path, AuthMiddle(), h.base.GuestRunAction(f))
	case "POST":
		guestRoute.POST(path, AuthMiddle(), h.base.GuestRunAction(f))
	case "PUT":
		guestRoute.PUT(path, AuthMiddle(), h.base.GuestRunAction(f))
	case "DELETE":
		guestRoute.DELETE(path, AuthMiddle(), h.base.GuestRunAction(f))
	default:
		panic(fmt.Sprintf(":%s method not allow", method))
	}
}
