FROM golang:1.22

WORKDIR /gateway

COPY . .

RUN go mod tidy

ENTRYPOINT [ "go run ./cmd/gateway-server" ]