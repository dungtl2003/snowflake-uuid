package insecure_test

import (
	"context"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	pb "dungtl2003/snowflake-uuid/internal/pb"
	"dungtl2003/snowflake-uuid/internal/server"
	"dungtl2003/snowflake-uuid/internal/tls"
)

const (
	bufSize = 1024 * 1024
)

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)

	idGeneratorServer, err := server.NewIdGeneratorServer(1, 1722613500)
	if err != nil {
		log.Fatalf("Failed to create ID generator server: %v", err)
	}

	go func() {
		if err := server.RunGrpcServer(idGeneratorServer, true, lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestGenerateId(t *testing.T) {
	ctx := context.Background()

	clientTlsCredentials, err := tls.LoadClientTlsCredentials()
	if err != nil {
		log.Fatalf("Failed to load client tls credentials: %v", err)
	}

	conn, err := grpc.NewClient(
		"passthrough:chatapp.com",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(clientTlsCredentials))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewIdGeneratorClient(conn)
	resp, err := client.GenerateId(ctx, &pb.GenerateIdRequest{})
	if err != nil {
		t.Fatalf("GenerateId failed: %v", err)
	}
	log.Printf("Response: %+v", resp)
}
