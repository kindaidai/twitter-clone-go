FROM golang:1.16

RUN mkdir /server
WORKDIR /server

COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

COPY *.go ./

EXPOSE 8080
