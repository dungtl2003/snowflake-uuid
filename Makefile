DOCKER_USERNAME ?= ilikeblue
DOCKER_FOLDER ?= ./docker
APPLICATION_NAME ?= snowflake-uuid
GIT_HASH ?= $(shell git log --format="%h" -n 1) 
CERT_CONFIG_ENV ?= dev
SERVER_PORT ?= 9000

_BUILD_ARGS_TAG ?= ${GIT_HASH}
_BUILD_ARGS_RELEASE_TAG ?= latest
_BUILD_ARGS_DOCKERFILE ?= Dockerfile

clean:
	$(info ==================== cleaning project ====================)
	rm -f internal/pb/*
	rm -f internal/data/x509/*.pem

gen:
	$(info ==================== generating new proto files ====================)
	protoc --proto_path=proto proto/*.proto  --go_out=:internal/pb --go-grpc_out=:internal/pb

cert:
	$(info ==================== generating new pem files ====================)
	cd internal/data/x509; ./cert_gen.sh --env ${CERT_CONFIG_ENV}; cd ../../..

test:
	$(info ==================== running tests ====================)
	go test ./...

benchmark:
	$(info ==================== running benchmark ====================)
	go test -bench=. ./...

_builder: clean gen cert test benchmark
	$(info ==================== building dockerfile ====================)
	docker buildx build --platform linux/amd64 --tag ${DOCKER_USERNAME}/${APPLICATION_NAME}:${_BUILD_ARGS_TAG} -f ${DOCKER_FOLDER}/${_BUILD_ARGS_DOCKERFILE} .

_pusher:
	docker push ${DOCKER_USERNAME}/${APPLICATION_NAME}:${_BUILD_ARGS_TAG}

_releaser:
	docker pull ${DOCKER_USERNAME}/${APPLICATION_NAME}:${_BUILD_ARGS_TAG}
	docker tag  ${DOCKER_USERNAME}/${APPLICATION_NAME}:${_BUILD_ARGS_TAG} ${DOCKER_USERNAME}/${APPLICATION_NAME}:${_BUILD_ARGS_RELEASE_TAG}
	docker push ${DOCKER_USERNAME}/${APPLICATION_NAME}:${_BUILD_ARGS_RELEASE_TAG}

_server:
	docker container run --rm --env-file .env ${DOCKER_USERNAME}/${APPLICATION_NAME}:${_BUILD_ARGS_TAG} -port=${SERVER_PORT}

_server_tls:
	docker container run --rm --env-file .env ${DOCKER_USERNAME}/${APPLICATION_NAME}:${_BUILD_ARGS_TAG} -port=${SERVER_PORT} -tls

run:
	./scripts/load_env.sh go run ./cmd/main.go -port=${SERVER_PORT}

run_tls:
	./scripts/load_env.sh go run ./cmd/main.go -port=${SERVER_PORT} -tls

server:
	$(MAKE) _server

server_tls:
	$(MAKE) _server_tls

server_%: 
	$(MAKE) _server \
		-e _BUILD_ARGS_TAG="$*-${GIT_HASH}" \
		-e _BUILD_ARGS_DOCKERFILE="Dockerfile.$*"

server_tls_%: 
	$(MAKE) _server_tls \
		-e _BUILD_ARGS_TAG="$*-${GIT_HASH}" \
		-e _BUILD_ARGS_DOCKERFILE="Dockerfile.$*"

build:
	$(MAKE) _builder
 
push:
	$(MAKE) _pusher
 
release:
	$(MAKE) _releaser

build_%: 
	$(MAKE) _builder \
		-e _BUILD_ARGS_TAG="$*-${GIT_HASH}" \
		-e _BUILD_ARGS_DOCKERFILE="Dockerfile.$*"
 
push_%:
	$(MAKE) _pusher \
		-e _BUILD_ARGS_TAG="$*-${GIT_HASH}"
 
release_%:
	$(MAKE) _releaser \
		-e _BUILD_ARGS_TAG="$*-${GIT_HASH}" \
		-e _BUILD_ARGS_RELEASE_TAG="$*-latest"

.PHONY:
	clean gen cert test benchmark server server_tls server_% server_tls_% \
	build run run_tls push release build_% push_% release_%
