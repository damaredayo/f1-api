FROM golang:1.17

WORKDIR /go/src/app
COPY . .

RUN go build -o f1-api

EXPOSE 8080

CMD ["./f1-api"]
