package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"dungtl2003/snowflake-uuid/internal/data"
	pb "dungtl2003/snowflake-uuid/internal/pb"
	"dungtl2003/snowflake-uuid/internal/server"
	"dungtl2003/snowflake-uuid/internal/tools"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	serverCertFile   = "x509/server_cert.pem"
	serverKeyFile    = "x509/server_key.pem"
	clientCACertFile = "x509/ca_cert.pem"
)

func loadTlsCredentials() (credentials.TransportCredentials, error) {
	pemClientCA, err := os.ReadFile(data.Path(clientCACertFile))
	if err != nil {
		return nil, err
	}

	certPools := x509.NewCertPool()
	if !certPools.AppendCertsFromPEM(pemClientCA) {
		return nil, fmt.Errorf("failed to add client CA's certificate")
	}

	cert, err := tls.LoadX509KeyPair(data.Path(serverCertFile), data.Path(serverKeyFile))
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

func runGrpcServer(idGeneratorServer *server.IdGeneratorServer, enableTls bool, lis net.Listener) error {
	var grpcServer *grpc.Server
	if !enableTls {
		grpcServer = grpc.NewServer()
	} else {
		tlsCredentials, err := loadTlsCredentials()
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

func main() {
	port := flag.Int("port", 0, "the server port")
	enableTLS := flag.Bool("tls", false, "enable SSL/TLS")
	flag.Parse()

	config := &tools.Configuration{}
	err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}

	idGeneratorServer, err := server.New(config.DatacenterId, config.WorkerId, config.Epoch, config.DatacenterIdBits, config.WorkerIdBits, config.SequenceBits)
	if err != nil {
		log.Fatalf("Failed to create ID generator server: %v", err)
	}

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen on port %v: %v", *port, err)
	}

	err = runGrpcServer(idGeneratorServer, *enableTLS, lis)
	if err != nil {
		log.Fatalf("Failed to serve grpc server on port %d: %v", *port, err)
	}
}
