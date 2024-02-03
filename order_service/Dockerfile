# Builder
FROM golang:1.21-alpine as builder

WORKDIR /app

RUN pwd && ls -l

COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN  go build -o /order_service github.com/mburaksoran/GetMobilCase/order_service

# Final docker image
FROM alpine:latest AS final

RUN apk update \
    && apk upgrade\
    && apk add --no-cache tzdata curl

#RUN apk --no-cache add bash DOCKER_BUILDKIT=0 docker build -t test -f Dockerfile .
ENV TZ Europe/Istanbul

WORKDIR /app
COPY --from=builder /order_service .
COPY --from=builder /app /app/

CMD [ "./order_service" ]