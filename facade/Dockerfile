FROM golang:1.22-alpine

RUN mkdir facade
COPY . /facade

WORKDIR /facade

RUN go build ./cmd/facade

CMD ["/facade/facade"]