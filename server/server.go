package server

import (
	"context"
	worker "dungtl2003/snowflake-uuid/internal"
	pb "dungtl2003/snowflake-uuid/proto"
	"log"
	"math/big"
)

type server struct {
	pb.UnimplementedIdGeneratorServer
	w *worker.Worker
}

func (s *server) GenerateId(ctx context.Context, in *pb.GenerateIdRequest) (*pb.GenerateIdResponse, error) {
	// log.Printf("Received request, context: %v", ctx)
	log.Printf("Received request")
	id, err := s.w.NextId()

	if err != nil {
		return nil, err
	}

	return &pb.GenerateIdResponse{Id: id.String()}, nil
}

func New(datacenterId *big.Int, workerId *big.Int, epoch *big.Int, datacenterIdBits *big.Int, workerIdBits *big.Int, sequenceBits *big.Int) *server {
	s := &server{w: worker.New(datacenterId, workerId, epoch, datacenterIdBits, workerIdBits, sequenceBits)}
	return s
}
