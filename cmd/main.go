package main

import (
	"log"
	"math/big"
	"net"

	pb "dungtl2003/snowflake-uuid/proto"
	s "dungtl2003/snowflake-uuid/server"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterIdGeneratorServer(grpcServer, s.New(big.NewInt(0), big.NewInt(1), nil, nil, nil, nil))
	log.Printf("Server is listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve grpc server on port 9000: %v", err)
	}
}
