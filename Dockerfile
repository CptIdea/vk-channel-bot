FROM golang:latest
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o build .cmd/main.go
CMD ["./build"]