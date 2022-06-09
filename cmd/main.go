package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"imman/crud_service/internal/service"
	"imman/crud_service/internal/service/repository"
	"imman/crud_service/protos/protos/crud_pb"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"google.golang.org/grpc"
)

var (
	port string
)

func main() {
	flag.StringVar(&port, "port", ":8002", "port")
	flag.Parse()

	config, err := pgx.ParseConfig("postgres://fprotimaru:1@localhost:5432/test_db?sslmode=disable")
	if err != nil {
		panic(err)
	}
	config.PreferSimpleProtocol = true

	sqldb := stdlib.OpenDB(*config)
	db := bun.NewDB(sqldb, pgdialect.New())

	repo := repository.NewPostRepository(db)

	uc := service.NewPostService(repo)

	server := grpc.NewServer()
	crud_pb.RegisterPostCRUDServiceServer(server, uc)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln(err)
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err = server.Serve(lis); err != nil {
			log.Fatalln(err)
		}
	}()

	log.Println("post_parser gRPC service is running on", port)
	<-quit
	log.Println("stopping post_parser gRPC service...")
	server.GracefulStop()
	log.Println("stopped post_parser gRPC service")
}
