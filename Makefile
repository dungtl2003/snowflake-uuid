gen:
	protoc --proto_path=proto proto/*.proto  --go_out=:internal/pb --go-grpc_out=:internal/pb

cert-dev:
	cd internal/data/x509; ./cert_gen.sh; cd ../../..

cert-prod:
	cd internal/data/x509; ./cert_gen.sh --env prod; cd ../../..

test:
	go test ./...

benchmark:
	go test -bench=. ./...

server:
	go run ./cmd/main.go -port=9000

server-tls:
	go run ./cmd/main.go -port=9000 -tls

clean:
	rm internal/pb/*

.PHONY:
	gen cert-dev cert-prod test benchmark server clean
