package insecure_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/test/bufconn"

	"dungtl2003/snowflake-uuid/internal/data"
	pb "dungtl2003/snowflake-uuid/internal/pb"
	"dungtl2003/snowflake-uuid/internal/server"
)

const (
	bufSize = 1024 * 1024

	serverCertFile   = "x509/server_cert.pem"
	serverKeyFile    = "x509/server_key.pem"
	serverCACertFile = "x509/ca_cert.pem"

	clientCertFile   = "x509/client_cert.pem"
	clientKeyFile    = "x509/client_key.pem"
	clientCACertFile = "x509/ca_cert.pem"
)

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)

	idGeneratorServer, err := server.New(1, 1722613500)
	if err != nil {
		log.Fatalf("Failed to create ID generator server: %v", err)
	}

	tlsCredentials, err := loadServerTlsCredentials()
	if err != nil {
		log.Fatalf("Failed to load server tls credentials: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(tlsCredentials),
	)

	pb.RegisterIdGeneratorServer(grpcServer, idGeneratorServer)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func loadServerTlsCredentials() (credentials.TransportCredentials, error) {
	return loadTlsCredentials(serverCACertFile, serverCertFile, serverKeyFile)
}

func loadClientTlsCredentials() (credentials.TransportCredentials, error) {
	return loadTlsCredentials(clientCACertFile, clientCertFile, clientKeyFile)
}

func loadTlsCredentials(caCertFile string, certFile string, keyFile string) (credentials.TransportCredentials, error) {
	pemCA, err := os.ReadFile(data.Path(caCertFile))
	if err != nil {
		return nil, err
	}

	certPools := x509.NewCertPool()
	if !certPools.AppendCertsFromPEM(pemCA) {
		return nil, fmt.Errorf("failed to add CA's certificate")
	}

	cert, err := tls.LoadX509KeyPair(data.Path(certFile), data.Path(keyFile))
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPools,
	}
	return credentials.NewTLS(config), nil
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestGenerateId(t *testing.T) {
	ctx := context.Background()

	clientTlsCredentials, err := loadClientTlsCredentials()
	if err != nil {
		log.Fatalf("Failed to load client tls credentials: %v", err)
	}

	conn, err := grpc.NewClient(
		"passthrough:whatever",
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
