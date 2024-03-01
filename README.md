Neiro Golang Developer Test Cases
====

We are creating a mobile app where people can chat with different characters (AI) and receive generated text/audio/video messages.
That’s why sometimes we need fast and easy-to-access storage to use with our mobile apps


# Test Assignment

Develop key-value storage (redis-like) that supports `Set`, `Get`, and `Delete` methods

## Features

- Use any approach to guarantee O(log(N)) or O(1) complexity to retrieve requested elements
- Expired data should be deleted automatically. TTL should be provided by an user, or any default value should be used. The operation should be implemented in the most effective way
- Don’t use the filesystem to store data (store everything in RAM)
- Don’t use external services, DBMS, or network data transfers within the algorithm

## Requirements

The code must be written in Go

You must provide a GitHub repo that includes:

- a `README.md` file that contains a description of your approach
- code that can be used to run the project, as well as the description and environment for it

## Service settings
```yaml
env: "prod"
version: "1.0.0"
ttl: 60
port: 8260
use_tracing: true
tracing_address: "http://host.docker.internal:14268/api/traces"
```

Service settings are stored in the `config/my_service/prod.yaml` file.
Values are
- `env` - environment (local/prod)
- `version` - service version
- `ttl` - default time to live for keys (seconds)
- `port` - service port
- `use_tracing` - use tracing or not
- `tracing_address` - tracing address

## Service run

```shell
export MY_SERVICE_CONFIG_PATH=config/my_service/prod.yaml && go run ./cmd/my_service
```

## Running services via docker containers

```shell
docker-compose up -d
```

## Stop services in docker containers

```shell
docker-compose down
```

# How to Test Using Postman

### Test GET request
```shell
GET http://localhost:8260/api/get/1
```
request body
```json

```
response 

200 OK

```json
{
    "key": "1",
    "value": "value1"
}
```

### Test SET request
```shell
POST http://localhost:8260/api/set
```
request body
```json
{
  "key": "1",
  "value": "value1"
}
```
response

200 OK

```json

```

### Test DELETE request
```shell
DELETE http://localhost:8260/api/delete/1
```
request body
```json

```
response

200 OK


### Added file compression with UPX

```shell  
upx --best my_service
upx --ultra-brute --no-lzma my_service
upx --ultra-brute --lzma my_service
```

| Parameter               | Result    | Origin    | Ratio |
|-------------------------|-----------|-----------|-------|
| --best                  | 3 670 032 | 8 196 674 |   44.77%    |
| --ultra-brute --no-lzma | 3 670 032 | 8 196 674 |  44.77%   |
| --ultra-brute --lzma    | 2 277 392 | 8 196 674 |  27.78%   |

