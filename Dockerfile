# build stage
FROM golang:1.18-alpine AS builder

RUN mkdir /app

WORKDIR /app

COPY . /app
RUN go build -v -o simplebank .
RUN chmod +x simplebank

RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz

# stage 2
FROM alpine:latest

RUN mkdir /app

WORKDIR /app

COPY --from=builder /app/simplebank /app
COPY --from=builder /app/migrate ./migrate

COPY db/migration ./db/migration
COPY app.env /app
COPY wait-for /app
COPY start.sh /app

CMD [ "/app/simplebank" ]
#ENTRYPOINT [ "/app/start.sh" ]