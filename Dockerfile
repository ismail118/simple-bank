# build stage
FROM golang:1.18-alpine AS builder

RUN mkdir /app

WORKDIR /app

COPY . /app
RUN go build -v -o simplebank .
RUN chmod +x simplebank

# stage 2
FROM alpine:latest

RUN mkdir /app

WORKDIR /app

COPY --from=builder /app/simplebank /app
COPY app.env /app
COPY wait-for /app
COPY db/migration /app/db/migration

CMD [ "/app/simplebank" ]