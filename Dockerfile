FROM golang:1.22.2-bullseye AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -x -o main main.go

FROM alpine:3.19.1

WORKDIR /app

COPY --from=builder /app/main .

COPY .env .
COPY token.json .
COPY credentials.json .

EXPOSE 8080

CMD [ "/app/main" ]