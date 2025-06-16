# syntax=docker/dockerfile:1

FROM golang:latest

WORKDIR /github.com/phthallo/guestbook

COPY docker_ip.sh ./

RUN chmod +x docker_ip.sh

ENTRYPOINT ["/docker_ip.sh"]

COPY go.mod go.sum ./

COPY main.go ./

COPY api ./api

COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /guestbook

CMD ["/guestbook"]

