FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go test ./... -v

RUN go build -o demoService .

FROM ubuntu:latest

WORKDIR /app
COPY --from=builder /app/demoService .

COPY .env .

COPY src/web/templates ./src/web/templates
COPY src/web/static ./src/web/static

CMD ["./demoService"]
