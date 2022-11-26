FROM golang:1.17 as builder
WORKDIR /app
COPY app /app
RUN go mod download &&  CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/main /app/main
CMD ["/app/main"]
