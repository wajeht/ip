FROM golang:1.26-alpine@sha256:bd14630652464086289693533d25b791aa9ae7481e784d7eac5d4c948e9736ea AS build

RUN apk --no-cache add curl wget

ENV ENV=production

WORKDIR /go/src/app

RUN wget -O GeoLite2-City.mmdb https://git.io/GeoLite2-City.mmdb

COPY . .

RUN go build -o ip .

EXPOSE 80

HEALTHCHECK CMD curl -f http://localhost:80/healthz || exit 1

CMD ["./ip"]
