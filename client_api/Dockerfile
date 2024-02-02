# Builder
FROM golang:1.21-alpine as builder

WORKDIR /app

RUN pwd && ls -l

COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN  go build -o /client_api github.com/mburaksoran/GetMobilCase/client_api

# Final docker image
FROM alpine:latest AS final

RUN apk update \
    && apk upgrade\
    && apk add --no-cache tzdata curl

#RUN apk --no-cache add bash DOCKER_BUILDKIT=0 docker build -t test -f Dockerfile .
ENV TZ Europe/Istanbul

WORKDIR /app
COPY --from=builder /client_api .
COPY --from=builder /app /app/

EXPOSE 8080
CMD [ "./client_api" ]