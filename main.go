package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"

	"egnite.app/microservices/user/config"
	"egnite.app/microservices/user/services/user"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
)

func startGRPC(wg *sync.WaitGroup) {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	userServer := user.Server{}

	user.RegisterUserServiceServer(grpcServer, &userServer)
	log.Println("gRPC server ready...")
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
	wg.Done()
}

func startHTTP(wg *sync.WaitGroup) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Connect to the GRPC server
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	// Register grpc-gateway
	rmux := runtime.NewServeMux()

	userClient := user.NewUserServiceClient(conn)
	err = user.RegisterUserServiceHandlerClient(ctx, rmux, userClient)
	if err != nil {
		log.Fatal(err)
	}
	handler := cors.Default().Handler(rmux)
	log.Println("rest server ready...")

	err = http.ListenAndServe(":8000", handler)
	if err != nil {
		log.Fatal(err)
	}
	wg.Done()

}

func main() {
	var wg sync.WaitGroup

	config.InitialiseEnvironment()
	wg.Add(1)
	go startGRPC(&wg)

	wg.Add(1)
	go startHTTP(&wg)

	wg.Wait()

}
