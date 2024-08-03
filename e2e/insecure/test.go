package insecure_test

import (
	"context"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	pb "dungtl2003/snowflake-uuid/internal/pb"
	"dungtl2003/snowflake-uuid/internal/server"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()

	idGeneratorServer, err := server.New(1, 1722613500)
	if err != nil {
		log.Fatalf("Failed to create ID generator server: %v", err)
	}

	pb.RegisterIdGeneratorServer(s, idGeneratorServer)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestGenerateId(t *testing.T) {
	ctx := context.Background()

	conn, err := grpc.NewClient(
		"passthrough:whatever",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
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
