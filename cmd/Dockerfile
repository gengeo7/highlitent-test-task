# syntax=docker/dockerfile:1

FROM golang:1.25-alpine

RUN mkdir /app
 
COPY . /app
 
WORKDIR /app
 
RUN go build -o server cmd/main.go 
 
CMD [ "/app/server" ]



