FROM golang:1.26-alpine@sha256:a6a091eac01ceac4b97496fe2957a49b6cdd83365337d5f46f6f73710424e805 AS build

RUN apk --no-cache add curl wget

ENV ENV=production

WORKDIR /go/src/app

RUN wget -O GeoLite2-City.mmdb https://git.io/GeoLite2-City.mmdb

COPY . .

RUN go build -o ip .

EXPOSE 80

HEALTHCHECK CMD curl -f http://localhost:80/healthz || exit 1

CMD ["./ip"]
