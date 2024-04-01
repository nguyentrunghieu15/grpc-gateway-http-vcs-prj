FROM golang:1.22

WORKDIR /gateway

COPY . .

RUN go mod tidy

CMD [ "go","run","./cmd/gateway-server" ]