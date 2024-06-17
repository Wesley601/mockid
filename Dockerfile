FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

RUN go install github.com/a-h/templ/cmd/templ@latest

RUN templ generate

# use -extldflags "-static" for go-sqlite3 building
RUN GOOS=linux GOARCH=amd64 go build -o mockid -ldflags='-s -w -extldflags "-static"' .

FROM scratch

ENV GO_ENV="PROD"

COPY --from=builder /app/mockid /app/mockid

ENTRYPOINT ["/app/mockid"]
