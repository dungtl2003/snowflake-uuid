package tls

import (
	"crypto/tls"
	"crypto/x509"
	"dungtl2003/snowflake-uuid/internal/data"
	"fmt"
	"os"

	"google.golang.org/grpc/credentials"
)

const (
	serverCertFile   = "x509/server_cert.pem"
	serverKeyFile    = "x509/server_key.pem"
	serverCACertFile = "x509/ca_cert.pem"

	clientCertFile   = "x509/client_cert.pem"
	clientKeyFile    = "x509/client_key.pem"
	clientCACertFile = "x509/ca_cert.pem"
)

func LoadClientTlsCredentials() (credentials.TransportCredentials, error) {
	pemCA, err := os.ReadFile(data.Path(serverCACertFile))
	if err != nil {
		return nil, err
	}

	certPools := x509.NewCertPool()
	if !certPools.AppendCertsFromPEM(pemCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	cert, err := tls.LoadX509KeyPair(data.Path(clientCertFile), data.Path(clientKeyFile))
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPools,
	}
	return credentials.NewTLS(config), nil
}

func LoadServerTlsCredentials() (credentials.TransportCredentials, error) {
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
