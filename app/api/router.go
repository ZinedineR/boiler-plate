package api

import (
	"boiler-plate/internal/base/handler"
	pb "boiler-plate/proto/helloworld/v1"
	pbTransaction "boiler-plate/proto/transaction/v1"
	pbUsers "boiler-plate/proto/users/v1"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func (h *HttpServe) setupGRPCUserRouter() {
	pb.RegisterServiceServer(h.GRPCUserServer, h.base.GRPCHandler)
	pbUsers.RegisterServiceServer(h.GRPCUserServer, h.UsersHandler.GRPCHandler)
}

func (h *HttpServe) setupGRPCUserGateway() {
	client, err := grpc.NewClient(
		"0.0.0.0:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	if err := pbUsers.RegisterServiceHandler(context.Background(), h.GRPCGateway, client); err != nil {
		log.Fatal("failed to register users gateway")
	}
	h.router.Group("/api/v2/*{grpc_gateway}").Any("", gin.WrapH(h.GRPCGateway))
}

func (h *HttpServe) setupGRPCTransactionRouter() {
	pb.RegisterServiceServer(h.GRPCTransactionServer, h.base.GRPCHandler)
	pbTransaction.RegisterServiceServer(h.GRPCTransactionServer, h.TransactionHandler.GRPCHandler)
}

func (h *HttpServe) setupGRPCTransactionGateway() {
	client, err := grpc.NewClient(
		"0.0.0.0:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	if err := pbTransaction.RegisterServiceHandler(context.Background(), h.GRPCGateway, client); err != nil {
		log.Fatal("failed to register transaction gateway")
	}
	//h.router.Group("/api/v2/*{grpc_gateway}").Any("", gin.WrapH(h.GRPCGateway))
}

func (h *HttpServe) setupUsersRouter() {
	h.GuestRoute("GET", "/users", h.UsersHandler.Find)
	h.GuestRoute("POST", "/users", h.UsersHandler.Create)
	h.GuestRoute("PUT", "/users/:id", h.UsersHandler.Update)
	h.GuestRoute("GET", "/users/:id", h.UsersHandler.Detail)
	h.GuestRoute("DELETE", "/users/:id", h.UsersHandler.Delete)
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

func (h *HttpServe) GRPCRoute(method, path string, f handler.HandlerFnInterface) {
	guestRoute := h.router.Group("/v1/*{grpc_gateway}")
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
