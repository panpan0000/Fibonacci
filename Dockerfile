FROM golang:1.8
WORKDIR /go/src/app
ADD . /go/src/app
RUN export GOPATH=`pwd` &&  go build -o fibonacci src/main.go
CMD ./fibonacci
EXPOSE 8008
