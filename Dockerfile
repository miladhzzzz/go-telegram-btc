FROM golang:1.21-alpine as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN cd cmd && go build -o bot-server

FROM alpine:latest as server

WORKDIR /app

COPY --from=builder /app/cmd/bot-server .

RUN chmod +x ./bot-server

ENV TELEGRAM_TOKEN=$TELEGRAM_TOKEN

CMD ["./bot-server"]
