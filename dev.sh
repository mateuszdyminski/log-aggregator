#!/bin/bash

usage() {
  cat <<EOF
Usage: $(basename $0) <command>

Wrappers around core binaries:
    frontend                Runs the frontend.
    backend                 Runs the backend.
    web                     Runs the web - default port 8001.
    nsqStart                Runs nsq dockers.
    nsqStop                 Stop nsq dockers.
    deps                    Installs all dependencies.
    pushFrontend            Builds and pushes container to quay.io. You have to be login first! type 'docker login quay.io'.
    pushBackend             Builds and pushes container to quay.io. You have to be login first! type 'docker login quay.io'.
    pushWeb                 Builds and pushes container to quay.io. You have to be login first! type 'docker login quay.io'.
    pushNsqd                Pushes container to quay.io. You have to be login first! type 'docker login quay.io'.
    pushLb                  Pushes container to quay.io. You have to be login first! type 'docker login quay.io'.
    pushAll                  Pushes all containers to quay.io.
    build                   Builds all GO components.
EOF
  exit 1
}

GO=${GO:-$(which go)}
export GOPATH="$PWD"
export PATH="$PATH:$GOPATH/bin"

VERSION=${IMT_VERSION-"$(git describe --always)"}

install_deps() {
  echo "Installing dependencies..."

  export GOPATH="$PWD/web"
  $GO get github.com/bitly/go-nsq
  $GO get github.com/mateuszdyminski/glog

  export GOPATH="$PWD/backend"
  $GO get github.com/bitly/go-nsq
  $GO get github.com/gorilla/websocket
  $GO get github.com/BurntSushi/toml

  export GOPATH="$PWD/frontend"
  $GO get github.com/gorilla/mux
}

buildFrontend() {
  echo "Building frontend..."
  set -e
  export GOOS="linux" && export GOARCH="amd64" && export GOPATH="$PWD/frontend" && go build -o frontend/bin/server server
  set +e
  echo "Frontend builed with success!"
}

buildBackend() {
  echo "Building backend..."
  set -e
  export GOOS="linux" && export GOARCH="amd64" && export GOPATH="$PWD/backend" && go build -o backend/bin/backend backend
  set +e
  echo "Backend builed with success!"
}

buildWeb() {
  echo "Building web..."
  set -e
  export GOOS="linux" && export GOARCH="amd64" && export GOPATH="$PWD/web" && go build -o web/bin/main emiter
  set +e
  echo "Web builed with success!"
}

stopNsq() {
  echo "Stop nsq dockers..."
  set -e
  docker ps | grep quay.io/mateuszdyminski/nsq | awk '{ print $1}' | xargs --no-run-if-empty docker kill
  set +e
  echo "Nsq stopped with success!"
}

startNsq() {
  echo "Starting nsq dockers..."
  set -e
  docker run -p 4150:4150 -p 4151:4151 -d quay.io/mateuszdyminski/nsq nsqd --broadcast-address=172.17.42.1 --lookupd-tcp-address=172.17.42.1:4160
  docker run -p 4160:4160 -p 4161:4161 -d quay.io/mateuszdyminski/nsq nsqlookupd --broadcast-address=172.17.42.1
  set +e
  echo "Nsq started with success!"
}

pushFrontend() {
  echo "Pushing frontend container started..."
  buildFrontend
  set -e
  docker build -t quay.io/mateuszdyminski/frontend:latest "$PWD/frontend/."
  docker push quay.io/mateuszdyminski/frontend
  set +e
  echo "Frontend pushed to quay.io!"
}

pushBackend() {
  echo "Pushing backend container started..."
  buildBackend
  set -e
  docker build -t quay.io/mateuszdyminski/backend:latest "$PWD/backend/."
  docker push quay.io/mateuszdyminski/backend
  set +e
  echo "Backend pushed to quay.io!"
}

pushWeb() {
  echo "Pushing web container started..."
  buildWeb
  set -e
  docker build -t quay.io/mateuszdyminski/web:latest "$PWD/web/."
  docker push quay.io/mateuszdyminski/web
  set +e
  echo "Web pushed to quay.io!"
}

pushNsqd() {
  echo "Pushing Nsqd container started..."
  set -e
  docker build -t quay.io/mateuszdyminski/nsq:latest "$PWD/nsq-confd/."
  docker push quay.io/mateuszdyminski/nsq
  set +e
  echo "Nsq pushed to quay.io!"
}

pushLb() {
  echo "Pushing nginx lb container started..."
  set -e
  docker build -t quay.io/mateuszdyminski/nginx_lb:latest "$PWD/nginx-lb/."
  docker push quay.io/mateuszdyminski/nginx_lb
  set +e
  echo "Nginx lb pushed to quay.io!"
}

CMD="$1"
shift
case "$CMD" in
  deps)
    install_deps
  ;;
  frontend)
    buildFrontend
    exec "$PWD/frontend/bin/server" --dir="$PWD/frontend/app" --p="9001"
  ;;
  backend)
    buildBackend
    echo "NsqlookupdAddresses = [\"127.0.0.1:4161\"]" > /tmp/backend.toml
    exec "$PWD/backend/bin/backend" --p="8090" --config="/tmp/backend.toml"
  ;;
  web)
    buildWeb
    exec "$PWD/web/bin/main" --logtostderr="1" --toChannel="true" --nsqdAddr="127.0.0.1:4150" --p="8001" --host="127.0.0.1"
  ;;
  build)
    buildFrontend
    buildBackend
    buildWeb
    echo "All builds done!"
  ;;
  nsqStart)
    startNsq
  ;;
  nsqStop)
    stopNsq
  ;;
  pushFrontend)
    pushFrontend
  ;;
  pushBackend)
    pushBackend
  ;;
  pushWeb)
    pushWeb
  ;;
  pushNsqd)
    pushNsqd
  ;;
  pushLb)
    pushLb
  ;;
  pushAll)
    pushFrontend
    pushBackend
    pushWeb
    pushNsqd
    pushLb
  ;;
  *)
    usage
  ;;
esac
