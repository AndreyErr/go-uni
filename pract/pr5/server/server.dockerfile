FROM golang:alpine AS builder
WORKDIR /src
COPY . .
RUN go build -o /src/main .

FROM alpine:latest
COPY --from=builder /src/main /app/main
CMD ["/app/main"]