FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bgp-manager .

FROM python:3.12-alpine

RUN pip install --no-cache-dir ansible

RUN apk add --no-cache docker-cli

WORKDIR /app

COPY --from=builder /app/bgp-manager .

EXPOSE 8080

CMD ["./bgp-manager"]
