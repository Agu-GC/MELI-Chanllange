FROM golang:1.24.2 AS builder

WORKDIR /app

COPY ./api/go.mod ./api/go.sum ./
RUN go mod download

COPY ./api/ .
COPY ./ssl ./ssl

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/ssl /app/ssl

EXPOSE 443

CMD ["./main"]