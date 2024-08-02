package server

import (
	"context"
	pb "dungtl2003/snowflake-uuid/internal/pb"
	"dungtl2003/snowflake-uuid/internal/worker"
	"log"
	"math/big"
)

type IdGeneratorServer struct {
	pb.UnimplementedIdGeneratorServer
	w *worker.Worker
}

func (s *IdGeneratorServer) GenerateId(ctx context.Context, in *pb.GenerateIdRequest) (*pb.GenerateIdResponse, error) {
	// log.Printf("Received request, context: %v", ctx)
	log.Printf("Received request")
	id, err := s.w.NextId()

	if err != nil {
		return nil, err
	}

	return &pb.GenerateIdResponse{Id: id.String()}, nil
}

func New(datacenterId *big.Int, workerId *big.Int, epoch *big.Int, datacenterIdBits *big.Int, workerIdBits *big.Int, sequenceBits *big.Int) (*IdGeneratorServer, error) {
	w, err := worker.New(datacenterId, workerId, epoch, datacenterIdBits, workerIdBits, sequenceBits)

	if err != nil {
		return nil, err
	}

	s := &IdGeneratorServer{w: w}
	return s, nil
}
