# build stage
FROM golang:1.9.2-alpine AS build-env
ADD . /src/go-schema-validator
ENV GOPATH=/
RUN apk add --no-cache git && cd /src/go-schema-validator && go get ./... && go build -o goapp

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /src/go-schema-validator/goapp /app/
ENTRYPOINT ./goapp