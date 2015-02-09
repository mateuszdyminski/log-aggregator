# WTF?


Log-aggregator shows how to aggregate logs from many diffents servers and show them in one place.   

Project contains 3 main parts:

- web - server generates logs which are pushed into kafka
- backendNodejs - backend server for log-aggregator - it pulls logs from kafka and send them to frontend via web-socket 
-  frontend - simple web app which displays logs in near realtime in the web browser  

This project was created at micro-hacathlon during web-socket workshop at Avaus Consulting(http://www.avaus.fi/). 

### Team:

- https://github.com/mateuszdyminski
- https://github.com/viru
- https://github.com/mmalczewski

# SETUP

## Requirements

- GOLANG see: https://golang.org/doc/install
- DOCKER see: https://docs.docker.com/installation/
- FIG see: http://www.fig.sh/install.html
- NODEJS see: http://nodejs.org/

## RUN LOCAL DOCKER CLUSTER

### CREATE EXECUTABLE HTTP SERVERS

Install log generator dependencies - first time only 

```
export GOPATH=$PWD/web && go get github.com/Shopify/sarama github.com/mateuszdyminski/glog
```

Install frontend dependencies - first time only

```
export GOPATH=$PWD/frontend && go get github.com/gorilla/mux
```

Build log generator:

```
export GOOS=linux && export GOARCH=amd64 && export GOPATH=$PWD/web && go build -o web/bin/main emiter
```

Build frontend server:

```
export GOOS=linux && export GOARCH=amd64 && export GOPATH=$PWD/frontend && go build -o frontend/bin/server server
```

Run all containers:

```
fig up -d
```

### TEST

1. Check if all 6 containers are up and running:

```
docker ps
```

2. Open log aggregator: http://127.0.0.1:9001
3. Open 2 tabs in web browser: http://127.0.0.1:8001 and http://127.0.0.1:8002
4. Http logs should appear in log-aggregator: http://127.0.0.1:9001 