FROM golang:1.10
WORKDIR /go/src/app
ADD . /go/src/app
RUN export GOPATH=`pwd` &&  go build -o fibonacci webservice
CMD ./fibonacci
EXPOSE 8008
