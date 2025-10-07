# Stage 1: The Builder Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy all the necessary source code first
# Docker will place them inside /app
COPY ./edc-service ./edc-service
COPY ./shared ./shared

# -- The Fix is Here --
# Now, change the working directory to be inside your service's code
WORKDIR /app/edc-service

# Download dependencies from within the correct directory
RUN go mod download

# Build the application. The '.' now correctly refers to /app/core-service
# The output binary is placed at /app/main
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# ---

# Stage 2: The Final Stage
FROM alpine:latest

WORKDIR /app

# Copy the compiled binary from the builder stage.
# The source path /app/main is absolute and easy to find.
COPY --from=builder /app/main .

EXPOSE 50051

CMD ["./main"]