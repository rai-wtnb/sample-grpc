package main

import (
	"log"
	"os"
	"os/signal"
	"reversi/gen/pb"
	"net"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"reversi/server/handler"
)

func main() {
	port := 50051
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("failed to listen")
	}

	server := grpc.NewServer()

	pb.RegisterMatchingServiceServer(
		server,
		handler.NewMatchingHandler(),
	)

	pb.RegisterGameServiceServer(
		server,
		handler.NewGameHandler(),
	)

	reflection.Register(server)

	go func(){
		log.Println("start gRPC server")
		server.Serve(lis)
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	server.GracefulStop()
}