FROM golang:1.25.1

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN go get github.com/notedit/rtmp

COPY . .

RUN go build -o /rtmp-ingest

EXPOSE 1935

CMD ["/rtmp-ingest"]
