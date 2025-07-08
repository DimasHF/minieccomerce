FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o main ./main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 3000

CMD ["./main"]