FROM golang:1.26-alpine@sha256:1fb7391fd54a953f15205f2cfe71ba48ad358c381d4e2efcd820bfca921cd6c6 AS build

RUN apk --no-cache add curl wget

ENV ENV=production

WORKDIR /go/src/app

RUN wget -O GeoLite2-City.mmdb https://git.io/GeoLite2-City.mmdb

COPY . .

RUN go build -o ip .

EXPOSE 80

HEALTHCHECK CMD curl -f http://localhost:80/healthz || exit 1

CMD ["./ip"]
