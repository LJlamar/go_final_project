FROM golang:1.26.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o todo-list main.go

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/todo-list .

COPY web ./web

RUN ls -l ./todo-list

CMD ["./todo-list"]