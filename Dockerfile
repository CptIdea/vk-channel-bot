FROM golang:latest
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o build ./main.go
ENVFILE ./.env
CMD ["./build"]