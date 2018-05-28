# What it is 
This is a repo to provide a Restful web service, to dump the N Fibonacci numbers sequence(json), the `N` is given via the GET operation of the restful API.

It's written in Go-lang, with unit-test.
the docker file and kubernetes deployment(with HPA) is also provided.


# API Usage

#### start server
1. Please start the server first(either running from source code or running from docker, refer to below sessions for more details)
NOTE: it's hardcoded in 8008 port( *it was designed to an available port, but it takes complexity to client side.. *)

#### curl the API
2. if the server starts well, do below( assuming you have curl installed)

Please find below for Restful API Request from Client

Request From localmachine
```
curl  localhost:8008/v1/fib?num=22
```

From remote (denote the $IP is the server IP address)
```
curl  $IP:8008/v1/fib?num=2222
```

the response data is in JSON format(an array), so you can use json tool to decode

example
```
curl  $IP:8008/v1/fib?num=4 | jq '.'
```
the output is
```
[
 0,
 1,
 1,
 2
]
```
**NOTE** of below concern cases:

* if the value of `num` is negative, the http response code will be 400 (bad request)
* if the value of `num` is 0, the reponse will be blank( `[]` in json)

# Source Code Walk Through
* src/
	* main.go : the http entrance point
	* main_test.go : the unit-test/function-test for http API
	* specs/
		* expected.go : a common header file for testing, listing the expected fibonacci sequences 
	* fib/
		* fib.go : Fibonacci caculator with DP algorithm (not recursive)
		* fib_test.go : unit-test of the fib.go
	* deploy/
		* fibnacci-deployment.yaml:  the deployment with 4 replica set
        * fibonacci-service.yaml:    the service using NodePort as a simple example
        * fibonacci-hpa.yaml :       HPA scale out when heavy loading encounted and at most 10 replicas


# Build/Run from Source Code

Prerequist : Go-Lang installed and $GOROOT Path setup (refer to Go-lang offical documents)

```
export GOPATH=`pwd`
go build -o fibonacci src/main.go
./fibonacci
```
NOTE: if `8008` port is occupied. it will complain
```
 listen tcp :8008: bind: address already in use
```


# Tests

### Restful API Server Test:

```
go test src/main.go src/main_test.go
```

### Test the sub-routine of fibonacci()
```
cd src/fib

# Unit Test
go test fib.go  fib_test.go  -v

# Bench Mark Test
go test fib.go  fib_test.go -test.bench=".*"
```


# Container/Kubernetest Deployment

### Docker build/run


##### docker BUILD:
and saving as image tag `fibonacci`
```
docker build ./ -t fibonacci
```

##### docker RUN:
run the images and bind the 8008 host port with 8008 port inside container 
```
docker run -d --rm --name Fibonacci -p 8008:8008 fibonacci
```

### Kubernetes Deployment

To achieve the simple load balance and scale out, Kubernetes is utilized here to provide a simple way to make it.
NOTE: QoS and other improvement T.B.D

```
kubectl apply -f  deploy/
```

* deployment: it will create 4 replica set of Fibonacci Rest API server.
* service. it will use NodePort to export the service
* HPA(horizontal-pod-autoscale): when in heavy loading, it's setup that once any POD is suffering from consuming with 50% CPU loading, at most 10 relicas in total will be expanding to scale out and reduce the loading to single pod. 


## To Do
