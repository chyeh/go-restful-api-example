FROM golang:1.8-alpine
RUN mkdir -p /go/src
ADD . /go/src/app/
WORKDIR /go/src/app
RUN go install .
CMD ["app"]
EXPOSE 80
