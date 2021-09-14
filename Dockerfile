# docker build -t igorwgitlab/cupcake-rdb:latest .
# docker run -it -v $(pwd)/redis-data:/mnt:ro igorwgitlab/cupcake-rdb:latest /mnt/dump.rdb

FROM golang:1.17-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN rm dump
RUN go build cmd/dump/dump.go

FROM alpine:latest
WORKDIR /
COPY --from=builder /app/dump ./
ENTRYPOINT ["./dump"]
