FROM golang:1.17-alpine

WORKDIR /work

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o get-docker-image

ENTRYPOINT ["./get-docker-image"]