FROM golang:1.21

WORKDIR /go/delivery


COPY ./ ./

RUN go build -o main .

CMD ["./main"]