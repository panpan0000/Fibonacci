* [What it is ](README.md#what-it-is)
* [Usage ](README.md#usage)
* [Source Code Walk Through](README.md#source-code-walk-through)
* [Build & Run from src](README.md#buildrun-from-source-code)
* [Tests](README.md#tests)
* [Docker & Kubernetes ](README.md#containerkubernetes-deployment)
* [Next Steps](README.md#to-do)



# What it is 
This is a repo to provide a Restful web service, to dump the N Fibonacci numbers sequence(json), the `N` is given via the GET operation of the restful API.

It's written in Go-lang, with unit-test.
the docker file and kubernetes deployment(with HPA) is also provided.


# Usage

#### 1. start server
Please start the server first
either 
* [running from source code](README.md#buildrun-from-source-code) or 
* [running from docker](README.md#docker-run), refer to below sessions for more details)

NOTE: service is hardcoded to bind 8008 port( *it was designed to an available port, but it takes complexity to client side.. *)

#### 2.curl the API
2. if the server starts well, do below( assuming you have curl installed)

Please find below for Restful API Request from Client

denote the `$IP` is the server IP address, `$N` is the number of the sequence.
```
curl  $IP:8008/v1/fib?num=$N
```
OR
```
curl  $IP:8008/v1/fibonacci?num=$N
```

the response data is in JSON format(an array), so you can use json tool to decode

example
```
curl  localhost:8008/v1/fib?num=4 | jq '.'
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
		* fib.go : Fibonacci caculator (not recursive)
		* fib_test.go : unit-test of the fib.go
	* deploy/
		* fibnacci-deployment.yaml:  the deployment with 4 replica set
        * fibonacci-service.yaml:    the service using NodePort as a simple example
        * fibonacci-hpa.yaml :       HPA scale out when heavy loading encounted and at most 10 replicas
    * Dockerfile: build the docker images
    * .travis.yaml : the Travis CI config file. for each PR/commit, it will do the automation test/deploy.(deploy is not included so far)

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


# Container/Kubernetes Deployment

### Docker build/run

The docker images is push to public docker hub `panpan0000/fibonacci`

##### docker BUILD:
if you want to re-build by yourself, do below and it will be saved as image tag `fibonacci`
```
docker build ./ -t fibonacci
```

##### docker RUN:
run the images and bind the 8008 host port with 8008 port inside container 
```
docker pull panpan0000/fibonacci
docker run -d --rm --name Fibonacci -p 8008:8008 panpan0000/fibonacci
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
It's kind of rush for this pilot project. There're some more things worthy as a production projects.
examples:

* **Extremely Large N**: if the N is extremely large:
   * if the output array is tooooo long, the http I/O will suffer from **timeout** (even server process crash as below). **Pagination** could be a solution for this case. adding `?page=$i` in the restful API, and limit the each "page" to a reasonable size of array.  
   * if this kind of request is so frequent, the duplication of caculation will waste a lot of CPU time. Similar as above, leveraging a **persistent** database is the way out.  
   * if it exceeds uint64 boundary, it would require different algorithm to do the caculation...
* Config files de-coupling. (currently the 8008 port is hard-coded)
* **Kubernetes optimization** , like QoS...
* **Heavy Workload**:  the kubernetes deployment now is for a modest workload, but for very high loading at short period of time(although I don't think this Fibonacci will be so popular..), but if so, the DP(Dynamic Processing) can be moved to distributed cloud cluster. In another word, the "cached" previous caculation results will be persisted. The simplest way is to put it into database. but a memory based KV cluster like redis will be a better solution though.
   
* I'm newbie to Go-lang(just days), there're lots of optimization oppotunity in the code.( why I chose GO instead of javascript+Node ? I don't know....)



-----------



example of server crash for tooo long http I/O due to very large N(100000000)
```
goroutine 51 [IO wait, 1 minutes]:
internal/poll.runtime_pollWait(0x7fb29832ee30, 0x72, 0xc420037e58)
        /home/rackhd/go/src/runtime/netpoll.go:173 +0x57
internal/poll.(*pollDesc).wait(0xc420450098, 0x72, 0xffffffffffffff00, 0x6e53c0, 0x7f1520)
        /home/rackhd/go/src/internal/poll/fd_poll_runtime.go:85 +0x9b
internal/poll.(*pollDesc).waitRead(0xc420450098, 0xc42008af00, 0x1, 0x1)
        /home/rackhd/go/src/internal/poll/fd_poll_runtime.go:90 +0x3d
internal/poll.(*FD).Read(0xc420450080, 0xc42008af41, 0x1, 0x1, 0x0, 0x0, 0x0)
        /home/rackhd/go/src/internal/poll/fd_unix.go:157 +0x17d
net.(*netFD).Read(0xc420450080, 0xc42008af41, 0x1, 0x1, 0x0, 0x0, 0x0)
        /home/rackhd/go/src/net/fd_unix.go:202 +0x4f
net.(*conn).Read(0xc42000e028, 0xc42008af41, 0x1, 0x1, 0x0, 0x0, 0x0)
        /home/rackhd/go/src/net/net.go:176 +0x6a
net/http.(*connReader).backgroundRead(0xc42008af30)
        /home/rackhd/go/src/net/http/server.go:668 +0x5a
created by net/http.(*connReader).startBackgroundRead
        /home/rackhd/go/src/net/http/server.go:664 +0xce
```
