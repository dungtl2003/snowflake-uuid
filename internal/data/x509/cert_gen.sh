#!/bin/bash

CONFIG_FILE="./openssl.cnf"

gen() {
    rm -f *.pem
    # Create the CA certificate
    openssl req -x509 \
      -newkey rsa:4096 \
      -nodes \
      -days 3650 \
      -keyout ca_key.pem \
      -out ca_cert.pem \
      -subj "/C=VN/ST=Ha Noi/L=Ha Noi/O=chatapp/CN=chatapp_ca" \
      -config "$CONFIG_FILE" \
      -extensions ca \
      -sha256

    # Generate a server private key
    openssl genrsa -out server_key.pem 4096

    # Create a Certificate Signing Request (CSR) for the server
    openssl req -new \
      -key server_key.pem \
      -out server_csr.pem \
      -subj "/C=VN/ST=Ha Noi/L=Ha Noi/O=mychatapp/CN=mychatapp.local" \
      -config "$CONFIG_FILE" \
      -reqexts server

    # Sign the server CSR with the CA key to generate the server certificate
    openssl x509 -req \
      -in server_csr.pem \
      -CAkey ca_key.pem \
      -CA ca_cert.pem \
      -days 3650 \
      -set_serial 1000 \
      -out server_cert.pem \
      -extfile "$CONFIG_FILE" \
      -extensions server \
      -sha256

    # Verify the server certificate
    openssl verify -verbose -CAfile ca_cert.pem server_cert.pem

    # Generate a client private key
    openssl genrsa -out client_key.pem 4096

    # Create a Certificate Signing Request (CSR) for the Client
    openssl req -new \
      -key client_key.pem \
      -out client_csr.pem \
      -subj "/C=VN/ST=Ha Noi/L=Ha Noi/O=mychatapp/CN=mychatapp.local" \
      -config "$CONFIG_FILE" \
      -reqexts client

    # Sign the Client CSR with the Client CA Key to Generate the Client Certificate
    openssl x509 -req \
      -in client_csr.pem \
      -CAkey ca_key.pem \
      -CA ca_cert.pem \
      -days 3650 \
      -set_serial 1000 \
      -out client_cert.pem \
      -extfile "$CONFIG_FILE" \
      -extensions client \
      -sha256

    # Verify the Client Certificate
    openssl verify -verbose -CAfile ca_cert.pem client_cert.pem

    # Cleanup
    rm server_csr.pem client_csr.pem
}

main() {
    gen
}

main "$@";
