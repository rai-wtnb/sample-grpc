package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
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

	server := grpc.NewServer()
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
