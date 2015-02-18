# WTF?

Log-aggregator shows how to aggregate logs from many diffents servers and show them in one place.   

Project contains 3 main parts:

- web - server generates logs which are pushed into kafka
- backend - backend server for log-aggregator - it pulls logs from kafka and send them to frontend via web-socket 
- frontend - simple web app which displays logs in near realtime in the web browser  

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

## RUN LOCAL DOCKER CLUSTER

### CREATE EXECUTABLE HTTP SERVERS

Install dependencies - first time only 

```
./dev.sh deps
```

Build all components

```
./dev.sh build
```

Run single component(not in docker container):

```
./dev.sh {frontend, backend, web}
```

Run all components in docker containers:

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
4. Http server logs should appear in log-aggregator: http://127.0.0.1:9001 

# TODO:

* Create proper fleet files in coreos/ - fix docker links to ambassador pattern

* Add Dockerfile with HAproxy based on etcd

* Add logs filtering
  * Add host<public ip> as key message send to NSQ - web
  * Add peridiacally asking about all web servers from etcd and send this info to frontend(over WebSocket) - backend
  * Add logs filtering over host ip - fronend 
