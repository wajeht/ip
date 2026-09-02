FROM golang:1.27-alpine@sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125 AS build

RUN apk --no-cache add curl wget

ENV ENV=production

WORKDIR /go/src/app

RUN wget -O GeoLite2-City.mmdb https://git.io/GeoLite2-City.mmdb

COPY . .

RUN go build -o ip .

EXPOSE 80

HEALTHCHECK CMD curl -f http://localhost:80/healthz || exit 1

CMD ["./ip"]
