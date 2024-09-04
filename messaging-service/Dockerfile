FROM golang:1.20-alpine as builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

COPY . .

ENV CGO_ENABLED=1
RUN go build -o /app/main ./cmd/main.go

FROM alpine:3.15

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /app/config/local.yaml ./config/local.yaml

CMD ["/app/main"]
