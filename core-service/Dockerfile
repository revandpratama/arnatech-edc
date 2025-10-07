
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY ./core-service ./core-service
COPY ./shared ./shared


WORKDIR /app/core-service

RUN go mod download


RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .


FROM alpine:latest

WORKDIR /app


COPY --from=builder /app/main .

EXPOSE 8101

CMD ["./main"]