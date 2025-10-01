FROM golang:1.24.5

COPY ./ /go/src/backend-test/
WORKDIR /go/src/backend-test/

RUN go get ./

RUN go build -o main .

CMD ["/go/src/backend-test/main"]

EXPOSE 1112