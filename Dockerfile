FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/app

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/server .
COPY config.yaml .

EXPOSE 50051

CMD ["./server"]
