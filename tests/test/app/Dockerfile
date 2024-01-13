FROM golang:1.21.4 as prod
WORKDIR /app
COPY . .
RUN go build -o /app/main
CMD ["/app/main"]
