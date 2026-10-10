FROM golang:1.27-alpine@sha256:f92b6ef800e499660581efdabdf25d9d817a9d124eaf900924f0504e7e27e12d AS build

RUN apk --no-cache add curl wget

ENV ENV=production

WORKDIR /go/src/app

RUN wget -O GeoLite2-City.mmdb https://git.io/GeoLite2-City.mmdb

COPY . .

RUN go build -o ip .

EXPOSE 80

HEALTHCHECK CMD curl -f http://localhost:80/healthz || exit 1

CMD ["./ip"]
