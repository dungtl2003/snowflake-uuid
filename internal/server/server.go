package server

import (
	"context"
	pb "dungtl2003/snowflake-uuid/internal/pb"
	"dungtl2003/snowflake-uuid/internal/worker"
)

type IdGeneratorServer struct {
	pb.UnimplementedIdGeneratorServer
	w *worker.Worker
}

func (s *IdGeneratorServer) GenerateId(ctx context.Context, in *pb.GenerateIdRequest) (*pb.GenerateIdResponse, error) {
	id, err := s.w.NextId()

	if err != nil {
		return nil, err
	}

	return &pb.GenerateIdResponse{Id: id}, nil
}

func New(workerId int64, epoch int64) (*IdGeneratorServer, error) {
	w, err := worker.New(workerId, epoch)

	if err != nil {
		return nil, err
	}

	s := &IdGeneratorServer{w: w}
	return s, nil
}
