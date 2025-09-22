FROM golang:1.25.1

WORKDIR /app

COPY go.mod ./
RUN go mod download
RUN go get github.com/notedit/rtmp

COPY . .

RUN go build -o /transcoder

EXPOSE 1935

CMD ["/transcoder"]
