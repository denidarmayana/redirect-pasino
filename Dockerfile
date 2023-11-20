FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM alpine:latest
EXPOSE 8080
RUN apk add ca-certificates ca-certificates-doc ca-certificates-bundle && update-ca-certificates
COPY --from=builder /app/main /main
CMD [ "/main" ]
