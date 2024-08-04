# About project

this is a chat app's service that supports generating twitter's snowflake IDs

# Table of content

[prerequisites](#-prerequisites)<br>
[setup](#-setup)<br>
[getting started](#-getting-started)<br>
[run tests](#-run-tests)<br>
[CI](#-ci)<br>
[deployment (comming soon)](#-deploy)<br>

## ⇁ Prerequisites

you must have go and docker installed<br>

## ⇁ Setup

first, clone this project<br>

you need to have `.env` file in root project, in the file you
need `key=value` each line. See list of required environment
variables [here](#-list-of-available-environment-variables):<br>

## ⇁ List of available environment variables

# Environment variables

| Variable | Required | Purpose | Default |
| -------- | -------- | ------- | ------- |
| WORKER_ID | YES | ID of the current worker (NOTE: do not create 2 workers with the same ID or you can have duplicated IDs)| NONE |
| EPOCH | NO | epoch time (start time) in milliseconds | 1704067200 (2024-01-01T00:00:00) |

For the full .env file example, check out [this template](./templates/.env.template)

## ⇁ Getting Started

first, you need to have `.env` file inside root project. See more
in [here](#-list-of-available-environment-variables)<br>

you can build and test project using this command:
```shell
make build
```
this will clean your project, add necessery files, test project then build docker image<br>
you can choose build options like this: `make build_dev` or `make build_prod`

there are 2 ways to run the server:

#### ⇁ Insecure

when using this, you don't have to do anything else. Just run:
```shell
make run
```

#### ⇁ TLS handshake

this project is using mutual TLS method, that means both client and server need to provide certificates to each other. First, you need to generate certificates for both client and server by using `make cert` (you don't need to do this if you have already ran `make build`). Next, copy `client_cert.pem`, `client_key.pem`, `ca_cert.pem` and `ca_key.pem` to your client application. Then you need to include these credentials when connecting to the server. You can change certificate configuration by adding `CERT_CONFIG_ENV={env}` to your `make cert` command (`make cert CERT_CONFIG_ENV=prod` for example)or run `make build_prod`<br>
to run server in TLS mode, use this command:
```shell
make run_tls
```
note that you need to make sure CA certificate of both client and server are the same<br>
you can run with docker using `make server` or `make server_tls`

## ⇁ Run tests

you can run all tests with:
```shell
make test
```

or you can benchmark ID generator function with:
```shell
make benchmark
```

## ⇁ CI

if you have [act](https://github.com/nektos/act), run `act` to fully test your project's workflow locally (automatically adding latest docker image to docker hub comming soon)

## ⇁ Deploy

comming soon
