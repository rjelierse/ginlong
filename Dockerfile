FROM golang:1.16-alpine AS build

RUN mkdir /app
WORKDIR /app

COPY . .

RUN go get -d -v ./...
RUN go install -v ./cmd/ginlong-receiver

FROM alpine:3.13

RUN apk add ca-certificates

COPY --from=build /go/bin/ginlong-receiver /usr/local/bin

CMD ["ginlong-receiver"]
