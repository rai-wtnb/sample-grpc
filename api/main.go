package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"

	"github.com/rai-wtnb/sample-grpc/gen/api"
	"github.com/rai-wtnb/sample-grpc/handler"
)

func main() {
	port := 50051
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	zapLogger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	grpc_zap.ReplaceGrpcLogger(zapLogger)

	server := grpc.NewServer(
		grpc.UnaryIntercepter(
			grpc_middleware.ChainUnaryServer(
				grpc_zap.UnaryServerInterceptor(zapLogger),
				grpc_auth.UnaryServerInterceptor(auth),
			),
		),
	)
	reflection.Register(server)

	api.RegisterPankakeBakerServiceServer(
		server,
		handler.NewBakerHandler(),
	)

	go func() {
		log.Printf("start grpc server port %v", port)
		server.Serve(lis)
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("stopping grpc server...")
	server.GracefulStop()
}

func auth(ctx context.Context) (context.Context, error) {
	if token, err := grpc_auth.AuthFromMD(ctx, "bearer"); err != nil {
		return nil, err
	}

	if token != "hi/mi/tsu" {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid bearer token")
	}

	return context.WithValue(ctx, "UserName", "God"), nil
}
