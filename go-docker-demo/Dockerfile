FROM golang:latest as golang

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN go get github.com/gin-contrib/cors
RUN go get github.com/gin-gonic/gin
RUN go get github.com/lib/pq

RUN go build -o /go-docker-demo

EXPOSE 8080

CMD ["/go-docker-demo"]