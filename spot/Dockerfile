FROM golang:1.22-alpine

RUN mkdir spot
COPY . /spot

WORKDIR /spot

RUN go build ./cmd/spot

CMD ["/spot/spot"]