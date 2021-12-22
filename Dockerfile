FROM golang:1.17-alpine

WORKDIR /work

#copy go.mod and go.sum to prepare environment
COPY go.mod ./
COPY go.sum ./
RUN go mod download

#copy data and main.go files to run the tool
ADD data ./data
COPY *.go ./

RUN go build -o get-docker-image

ENTRYPOINT ["./get-docker-image"]