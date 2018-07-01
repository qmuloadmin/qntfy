FROM golang:1.9

WORKDIR /go/src/github.com/qmuloadmin/qntfy
COPY . /go/src/github.com/qmuloadmin/qntfy
COPY ./stats /go/src/github.com/qmuloadmin/qntfy/stats

RUN go get github.com/montanaflynn/stats
RUN go get github.com/gin-gonic/gin
RUN go build cli.go
RUN go build web.go

EXPOSE 8080

CMD ./web
