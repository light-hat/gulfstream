FROM golang:1.25.1

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /hls-origin

EXPOSE 1935

CMD ["/hls-origin"]
