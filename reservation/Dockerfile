FROM golang:1.22-alpine

RUN mkdir reservation
COPY . /reservation

WORKDIR /reservation

RUN go build ./cmd/reservation

CMD ["/reservation/reservation"]