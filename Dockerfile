FROM golang

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .


EXPOSE 9000
CMD CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go run avito_chat.go
