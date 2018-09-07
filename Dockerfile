FROM golang:1.8-alpine
RUN apk --update add postgresql-client && rm -rf /var/cache/apk/*
RUN mkdir -p /go/src
ADD . /go/src/app/
WORKDIR /go/src/app
RUN go install .
EXPOSE 80
