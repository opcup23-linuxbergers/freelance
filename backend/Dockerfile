FROM golang:alpine AS build

WORKDIR /src

COPY . /src/

RUN \
    go get -u && \
    go build -o /out

FROM alpine:latest

WORKDIR /app
COPY --from=build /out /app/out

EXPOSE 18000

ENTRYPOINT [ "/app/out" ]
