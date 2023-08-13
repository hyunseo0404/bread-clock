FROM golang:1.20-alpine AS build-stage
LABEL maintainer="hyunseo0404@gmail.com"

RUN apk add build-base

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir bin
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/breadclock ./cmd/breadclock/...

FROM alpine:latest AS release-stage
LABEL maintainer="hyunseo0404@gmail.com"

WORKDIR /breadclock

COPY --from=build-stage /app/bin/breadclock ./bin/breadclock
COPY --from=build-stage /app/config.env ./config.env
COPY --from=build-stage /app/docs ./docs

EXPOSE 80

RUN chmod -R 777 /breadclock/bin/*

ENTRYPOINT ["./bin/breadclock"]
