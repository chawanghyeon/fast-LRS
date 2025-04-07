FROM golang:1.20

WORKDIR /app

COPY main.go .

RUN go mod init app && \
    go get github.com/lib/pq github.com/google/uuid && \
    go build -o app

CMD ["./app"]
