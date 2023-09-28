FROM golang:1.21.1

ENV GO111MODULE=auto \
    CGO_ENABLED=1 \
    GOOS=linux

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o main .

CMD ["/build/main"]
