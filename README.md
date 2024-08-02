# About project

this is a chat app's service that supports generating twitter's snowflake IDs

# Table of content

[prerequisites](#-prerequisites)<br>
[setup](#-setup)<br>
[getting started](#-getting-started)<br>
[run tests](#-run-tests)<br>
[deployment (comming soon)](#-deploy)<br>

## ⇁ Prerequisites

you must have go installed<br>

## ⇁ Setup

first, clone this project<br>

you need to have `.env` file in root project, in the file you
need `key=value` each line. See list of required environment
variables [here](#-list-of-available-environment-variables):<br>

## ⇁ List of available environment variables

# Environment variables
DATACENTER_ID=0
WORKER_ID=1

| Variable | Required | Purpose | Default |
| -------- | -------- | ------- | ------- |
| DATACENTER_ID | YES | ID of the current datacenter | NONE |
| WORKER_ID | YES | ID of the current worker (NOTE: do not create 2 workers with the same ID with the same datacenter ID or you can have duplicated IDs)| NONE |
| EPOCH | NO | epoch time (start time) in milliseconds | 1704067200 (2024-01-01T00:00:00) |
| DATACENTER_ID_BITS | NO | max number of bits for datacenter ID | 5. You can use up to 32 different datacenter IDs |
| WORKER_ID_BITS | NO | max number of bits for worker ID | 5. You can use up to 32 different worker IDs |
| SEQUENCE_BITS | NO | max number of bits for sequence number | 12. There can be up to 4096 different IDs at the same time |

For the full .env file example, check out [this template](./templates/.env.template)

## ⇁ Getting Started

### ⇁ Development

first, you need to have `.env` file inside root project. See more
in [here](#-list-of-available-environment-variables)<br>

you can build and test project using this command:
```shell
make build-dev
```

there are 2 ways to run the server:

#### ⇁ Insecure

when using this, you don't have to do anything else. Just run:
```shell
make server
```

#### ⇁ TLS handshake

this project is using mutual TLS method, that means both client and server need to provide certificates to each other. First, you need to generate certificates for both client and server by using `make cert-dev` (you don't need to do this if you have already ran `make build-dev`). Next, copy `client_cert.pem`, `client_key.pem`, `ca_cert.pem` and `ca_key.pem` to your client application. Then you need to include these credentials when connecting to the server. To run server in TLS mode, use this command:
```shell
make server-tls
```
note that you need to make sure CA certificate of both client and server are the same


### ⇁ Production

comming soon

## ⇁ Run tests

you can run all tests with:
```shell
make test
```

or you can benchmark ID generator function with:
```shell
make benchmark
```

## ⇁ Deploy

comming soon
