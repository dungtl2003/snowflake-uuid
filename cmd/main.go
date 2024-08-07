package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"dungtl2003/snowflake-uuid/internal/server"
	"dungtl2003/snowflake-uuid/internal/tools"
)

func main() {
	port := flag.Int("port", 0, "the server port")
	enableTLS := flag.Bool("tls", false, "enable SSL/TLS")
	flag.Parse()

	if *enableTLS {
		log.Println("Setting up GRPC server with TLS handshake")
	} else {
		log.Println("Setting up GRPC server without TLS handshake")
	}

	log.Println("Loading configurations")
	config := &tools.Configuration{}
	err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}

	log.Println("Creating ID generator server")
	idGeneratorServer, err := server.NewIdGeneratorServer(config.WorkerId, config.Epoch)
	if err != nil {
		log.Fatalf("Failed to create ID generator server: %v", err)
	}

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen on port %v: %v", *port, err)
	}

	log.Printf("Serving GRPC server on port %d", *port)
	err = server.RunGrpcServer(idGeneratorServer, *enableTLS, lis)
	if err != nil {
		log.Fatalf("Failed to serve grpc server on port %d: %v", *port, err)
	}
}
