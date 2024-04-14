FROM golang:1.22.2-bullseye AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main main.go

FROM alpine:3.19.1

WORKDIR /app

COPY --from=builder /app/main .

COPY .env .
COPY token.json .
COPY credentials.json .

EXPOSE 8001

RUN chmod +x /app/main

RUN chmod +r /app/.env
RUN chmod +r /app/token.json
RUN chmod +r /app/credentials.json

RUN ls -la
RUN ls -la /app
RUN ls -la /app/main


CMD [ "/app/main" ]