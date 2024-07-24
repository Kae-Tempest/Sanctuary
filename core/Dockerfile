FROM golang:1.22-alpine AS builder

WORKDIR  /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=mod -o sanctuary-api

FROM scratch

WORKDIR /app
COPY --from=builder /app/sanctuary-api .

CMD ["./sanctuary-api"]