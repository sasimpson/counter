FROM golang:alpine as builder
WORKDIR /go/src/github.com/sasimpson/counter
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o service

FROM alpine:latest
WORKDIR /app/
COPY --from=builder /go/src/github.com/sasimpson/counter/service .
CMD ["./service"]
EXPOSE 5000
