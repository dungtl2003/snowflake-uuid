package server

import (
	pb "dungtl2003/snowflake-uuid/internal/pb"
	"dungtl2003/snowflake-uuid/internal/tls"
	"net"

	"google.golang.org/grpc"
)

func RunGrpcServer(idGeneratorServer *IdGeneratorServer, enableTls bool, lis net.Listener) error {
	var grpcServer *grpc.Server
	if !enableTls {
		grpcServer = grpc.NewServer()
	} else {
		tlsCredentials, err := tls.LoadServerTlsCredentials()
		if err != nil {
			return err
		}

		grpcServer = grpc.NewServer(
			grpc.Creds(tlsCredentials),
		)
	}

	pb.RegisterIdGeneratorServer(grpcServer, idGeneratorServer)

	return grpcServer.Serve(lis)
}
