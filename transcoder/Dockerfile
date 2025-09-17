FROM golang:1.19

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o /transcoder

EXPOSE 1935

CMD ["/transcoder"]
