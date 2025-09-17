FROM golang:1.19

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o /hls-origin

EXPOSE 1935

CMD ["/hls-origin"]
