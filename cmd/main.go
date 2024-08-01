package main

import (
	"log"
	"net"
	"os"
	"strconv"

	"dungtl2003/snowflake-uuid/internal/server"
	"dungtl2003/snowflake-uuid/internal/tools"
	pb "dungtl2003/snowflake-uuid/proto"

	"google.golang.org/grpc"
)

func main() {
	environment := os.Getenv("ENV")
	if environment == "" {
		environment = "development"
	}
	config := &tools.Configuration{}
	err := config.Load(environment)
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(config.Port))
	if err != nil {
		log.Fatalf("Failed to listen on port %v: %v", config.Port, err)
	}

	idGeneratorServer, err := server.New(config.DatacenterId, config.WorkerId, config.Epoch, config.DatacenterIdBits, config.WorkerIdBits, config.SequenceBits)
	if err != nil {
		log.Fatalf("Failed to create ID generator server: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterIdGeneratorServer(grpcServer, idGeneratorServer)
	log.Printf("Server is listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve grpc server on port %v: %v", strconv.Itoa(config.Port), err)
	}
}
