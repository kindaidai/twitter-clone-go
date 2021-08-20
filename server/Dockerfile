FROM golang:1.17

EXPOSE 8080

RUN mkdir /server
WORKDIR /server
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download
COPY . .
CMD go run main.go

