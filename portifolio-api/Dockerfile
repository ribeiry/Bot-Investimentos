FROM golang:1.23-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /src/main .

EXPOSE 8080

CMD ["./main"]
