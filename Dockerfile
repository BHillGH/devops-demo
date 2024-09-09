# Stage build
FROM golang:1.23-alpine AS build-stage

WORKDIR /build

COPY go.mod ./
RUN go mod download

COPY *.go ./
COPY /data/ ./data/

RUN CGO_ENABLED=0 GOOS=linux go build -o ./devops-demo

# Stage deploy
FROM alpine:latest

RUN apk update

WORKDIR /app

COPY --from=build-stage /build/devops-demo .
COPY --from=build-stage /build/data ./data

ENTRYPOINT ["./devops-demo"]