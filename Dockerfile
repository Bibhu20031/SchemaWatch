FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o schemawatch ./cmd

FROM alpine:3.19

WORKDIR /app

RUN adduser -D appuser
USER appuser

COPY --from=builder /app/schemawatch .

EXPOSE 8080

CMD ["./schemawatch"]
