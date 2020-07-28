FROM golang:alpine as builder

RUN apk update && apk upgrade && apk add --update alpine-sdk && \
    apk --no-cache add bash git make

WORKDIR /app/datahow-service

COPY . .

RUN make

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app

COPY --from=builder /app/datahow-service/dist/datahow .

CMD ["./datahow"]
