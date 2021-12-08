FROM golang:1.16-alpine AS build

WORKDIR /go/src/app

COPY . .

RUN go get -d ./...
RUN go install ./cmd/ginlong-proxy

FROM alpine:3.15

RUN apk add ca-certificates

COPY --from=build /go/bin/ginlong-proxy /usr/bin/

CMD [ "ginlong-proxy" ]
