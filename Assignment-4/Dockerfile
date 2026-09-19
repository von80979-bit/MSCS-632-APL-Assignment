# One image with both employee scheduler programs. Run one with `docker run -it --rm <image> scheduler-rust`
# or `scheduler-go`, and pass `-e MAX_EMPLOYEES_PER_SLOT=<n>` to change the slot maximum.

# Alpine builds Rust against musl, so the binary is static and runs on the small final image.
FROM rust:1-alpine AS rust-build
WORKDIR /build
COPY rust/ .
RUN cargo build --release --locked

# Without cgo the Go binary is static as well.
FROM golang:1.24-alpine AS go-build
WORKDIR /build
COPY go/ .
RUN CGO_ENABLED=0 go build -o scheduler-go .

FROM alpine:3
COPY --from=rust-build /build/target/release/scheduler-rust /usr/local/bin/scheduler-rust
COPY --from=go-build /build/scheduler-go /usr/local/bin/scheduler-go
