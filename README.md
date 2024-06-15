# Opentelemetry-Go-Distributed-Tracing

## Requirements

* Go (version >= 1.16)
* MySQL 8
* serve

## Install dependencies run below command

* For installation of `serve` see: [https://www.npmjs.com/package/serve](https://www.npmjs.com/package/serve)

```bash
npm install -g serve
```

## `.env` file setup

 * Copy the `.env.example` file to `.env` in the `root` directory.
 * Copy the `.env.example` file to `.env` in the `frontend` directory.


## Navigate `order/server.go` , `payment/server.go` & `users/server.go`.

```go
const serviceName = "<YOUR_SERVICE_NAME>"
```

## Start individual microservices using below commands

1) User Service

```sh
go run ./customers
```

2) Payment Service

```sh
go run ./payment
```

3) Order Service

```sh
go run ./order
```

4) Start the frontend using following command

```sh
serve -l 5000 frontend
```

5) Access 

```bash
http://localhost:5000/
http://localhost:8080/
http://localhost:8081/
http://localhost:8082/
```

## OTel Setup 

Install Otelcol-contribute [using this link](https://github.com/open-telemetry/opentelemetry-collector-releases/releases)


## Atatus collector Configuration

* you can use collector configuration file `atatus-collector.yaml` for send OTel data to Atatus.


## Run otel-contrib

```bash
./otelcol-contrib --config=<Your-Local-path>/atatus-collector.yaml
```