FROM rust:1.83-alpine AS builder

WORKDIR /app

COPY . .
RUN cargo build --release

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/target/release/* /app/server

EXPOSE 3000

CMD ["/app/server"]
