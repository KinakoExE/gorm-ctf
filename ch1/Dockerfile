FROM golang:1.17 as builder

RUN apt-get update && apt-get install build-essential -y

ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    GO111MODULE=on

WORKDIR /go/src/ch1
COPY ./ch1/app ./
RUN go mod download

RUN go build -a -ldflags '-linkmode external -extldflags "-static"' 

CMD ["go", "run", "main.go"]
