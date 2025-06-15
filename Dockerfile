# syntax=docker/dockerfile:1

FROM golang:latest

WORKDIR /github.com/phthallo/guestbook

COPY go.mod go.sum ./

COPY main.go ./

COPY api ./api

COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /guestbook

RUN ssh-keygen -t ed25519 -f id_ed25519

CMD ["/guestbook"]

