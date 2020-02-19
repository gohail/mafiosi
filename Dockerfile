FROM golang:latest AS builder
WORKDIR /mafiosi
ADD . .

# build binary
# disabled cgo for running in ansible
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GO111MODULE=on \
    go build \
        -a \
        -installsuffix cgo \
        -tags netgo \
        -o ./bin/app \
            ./main.go

FROM alpine:latest
WORKDIR /ipam
COPY --from=builder /mafiosi/bin/app .

EXPOSE 8080
# Use -host flag if needed "./app -host=yourhost:port" default:"localhost:8080"
CMD ["./app"]
