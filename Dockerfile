# Build Phase
FROM golang:alpine AS builder

ARG VERSION=1.0.0
ENV VERSION=$VERSION

RUN apk update && apk add gcc make openssl-libs-static zlib-static zstd-libs libsasl lz4-dev lz4-static zstd-static libc-dev musl-dev upx

WORKDIR /app
COPY . /app
ENV GO111MODULE=on
ENV my_service_CONFIG_PATH=config/my_service/prod.yaml
RUN make build_service
# compress binary
RUN upx --ultra-brute --lzma my_service

# Execution Phase
FROM alpine:latest

RUN apk --no-cache add ca-certificates \
	&& addgroup -S app \
	&& adduser -S app -G app

WORKDIR /app
# COPY --from=builder /app .
COPY --from=builder /app/my_service /app/my_service
COPY --from=builder /app/config/my_service/prod.yaml /app/config/my_service/prod.yaml
RUN chmod -R 777 /app
USER app

# Expose port to the outside world
EXPOSE 8260

# Command to run the executable
CMD ["./my_service"]
