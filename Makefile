GO111MODULE := auto
export GO111MODULE

lint:
	golangci-lint run ./...

test:
	go test -count=1 -race ./...

build_service:
	go build -tags musl -ldflags="-w -extldflags '-static' -X 'main.Version=$(VERSION)'" -o my_service neiro/cmd/my_service

check-swagger:
	which swagger

swagger: check-swagger
	GO111MODULE=on go mod vendor && GO111MODULE=off swagger generate spec -o ./doc/swagger.json --scan-models

serve-swagger: check-swagger
	swagger serve -F=swagger ./doc/swagger.json
