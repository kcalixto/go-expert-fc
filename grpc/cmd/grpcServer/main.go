package main

import (
	"database/sql"
	"fmt"
	"net"

	"github.com/kcalixto/go-expert-fc/grpc/internal/database"
	"github.com/kcalixto/go-expert-fc/grpc/internal/pb"
	"github.com/kcalixto/go-expert-fc/grpc/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create a new gRPC server
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	// Register the service
	pb.RegisterCategoryServiceServer(grpcServer, service.NewCategoryService(database.NewCategory(db)))

	// Listen and serve
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	fmt.Println("Server is running on port :50051")
	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
