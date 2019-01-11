FROM golang:latest
WORKDIR $GOPATH/src/github.com/liuyh73/dailyhub.service
ADD . $GOPATH/src/github.com/liuyh73/dailyhub.service
RUN go get github.com/liuyh73/dailyhub.service
RUN go build .
EXPOSE 9090
ENTRYPOINT ["./dailyhub.service"]