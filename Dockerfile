# syntax=docker/dockerfile:1
FROM golang:1.23.0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init


RUN go build -o app ./main.go

EXPOSE 8080

CMD ["./app"]
