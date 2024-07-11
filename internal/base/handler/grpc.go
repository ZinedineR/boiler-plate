package handler

import (
	hello "boiler-plate/proto/helloworld/v1"
	"context"
)

type GRPCHandler struct {
	hello.UnimplementedServiceServer
}

//func NewGRPCHandler(helloWorld pb.UnimplementedGreeterServer) *GRPCHandler {
//	return &GRPCHandler{helloWorld}
//}

// SayHello implements helloworld.GreeterServer
func (s *GRPCHandler) SayHello(ctx context.Context, in *hello.SayHelloRequest) (*hello.SayHelloResponse, error) {
	return &hello.SayHelloResponse{Message: "Hello " + in.Name}, nil
}
