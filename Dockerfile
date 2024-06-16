FROM golang:1.22-bookworm AS builder

WORKDIR /app

COPY . .

RUN go build -ldflags="-s -w" .

FROM scratch

COPY --from=builder /app/mockid /app/mockid

ENTRYPOINT ["/app/mockid"]
