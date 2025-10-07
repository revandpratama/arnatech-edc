
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY ./edc-service ./edc-service
COPY ./shared/ ../shared/
COPY ./api/ ../api/


WORKDIR /app/edc-service


RUN go mod download


RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /api ./api

EXPOSE 50051

CMD ["./main"]