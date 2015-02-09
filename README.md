# Requirements

- GOLANG
- DOCKER
- FIG
- NODEJS



# RUN LOCAL DOCKER CLUSTER


## CREATE EXECUTABLE HTTP SERVERS
```
Install log generator dependencies - first time only 
export GOPATH=$PWD/web && go get github.com/Shopify/sarama github.com/mateuszdyminski/glog

Install frontend dependencies - first time only
export GOPATH=$PWD/frontend && go get github.com/gorilla/mux

Build log generator:
export GOPATH=$PWD/web && go build -o web/bin/main emiter

Build frontend server:
export GOPATH=$PWD/frontend && go build -o frontend/bin/server server

Run all containers:
fig up -d
```

### CREATE WEB DOCKER IMAGE 
export GOPATH=$PWD && go build -o bin/main emiter && sudo docker build -t mateuszdyminski/log-aggregator .

### CREATE BACKEND NODE DOCKER IMAGE 
docker build -t mateuszdyminski/log-aggregator-node-backend backed-nodejs

### CREATE FRONTEND NODE DOCKER IMAGE 
docker build -t mateuszdyminski/log-aggregator-frontend frontend





### RUN DOCKER WITH WEB SERVER
docker run -p 127.0.0.1:8001:8001 --link logaggregator_kafka_1:kafka mateuszdyminski/log-aggregator

### RUN DOCKER WITH BACKEND NODE
docker run --link logaggregator_zookeeper_1:zk mateuszdyminski/log-aggregator-node-backend

### RUN DOCKER WITH FRONTEND  
docker run -p 127.0.0.1:9001:9001 --link TBD mateuszdyminski/log-aggregator-frontend
