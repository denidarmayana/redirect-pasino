# Stage 1: Build the Go application
FROM golang:alpine AS builder
WORKDIR /app
COPY . . 
RUN go build -o main main.go

# Stage 2: Create a minimal image with the built application
FROM alpine:latest
EXPOSE 8080
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/main /main
CMD [ "/main" ]

