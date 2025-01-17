FROM golang:1.22-alpine

RUN mkdir user
COPY . /user

WORKDIR /user

RUN go build ./cmd/user

CMD ["/user/user"]