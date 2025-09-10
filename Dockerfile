FROM golang:1.25.1 AS builder
WORKDIR /build
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o /build/spallet ./cmd/spallet && chmod +x /build/spallet

FROM alpine:3.22.1
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
WORKDIR app
COPY --from=builder /build/spallet /app/spallet
COPY ./assets/docs /app/assets/docs
CMD ["./spallet"]