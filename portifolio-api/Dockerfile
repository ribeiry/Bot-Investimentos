FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/api/

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /src/main .

EXPOSE 8080

CMD ["./main"]