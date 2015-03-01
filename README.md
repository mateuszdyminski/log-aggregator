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
- ANSIBLE see: http://docs.ansible.com/intro_installation.html

## RUN LOCAL DOCKER CLUSTER

### Prepare all containers

Install dependencies - first time only 

```
./dev.sh deps
```

Build all components

```
./dev.sh build
```

### Run local cluster

Run all components as docker containers:

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

## RUN CORE OS CLUSTER - VAGRANT

Go to infra directory

```
cd infra
```

Go to https://discovery.etcd.io/new

```
https://discovery.etcd.io/new
```

Copy token

```
https://discovery.etcd.io/04bf190ab730ef19c53241a27644dd4a <- this is token
```

Paste it in infra/vagrant/user-data <token> tag

Launch vagrant

```
vagrant up
```

When all nodes are up and running go to:

```
cd infra/ansible
```

We have to edit private key for vagrant first, open file ssh.config and change following line:

```
IdentityFile /home/md/.vagrant.d/insecure_private_key => IdentityFile /path/to/your/vagrant/private/key
```

To prepare coreOS cluster for ansible run:

```
ansible-playbook -i hosts/vagrant-inventory tasks/roles/bootstrap.yml --tags=destroy
```

Verify if eveything works fine:

```
ansible -i hosts/vagrant-inventory coreos -m ping
```

To run all services in cluster:

```
ansible-playbook -i hosts/vagrant-inventory tasks/roles/services.yml
```

To stop all services in cluster:

```
ansible-playbook -i hosts/vagrant-inventory tasks/roles/services.yml --tags=destroy
```

# TODO:

* Add logs filtering
  * Add peridiacally asking about all web servers from etcd and send this info to frontend(over WebSocket) - backend
  * Add logs filtering over host ip - fronend 
