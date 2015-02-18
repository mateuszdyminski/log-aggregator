#!/bin/bash

usage() {
  cat <<EOF
Usage: $(basename $0) <command>

Wrappers around core binaries:
    frontend               Runs the frontend.
    backend                Runs the backend.
    web                    Runs the web - default port 8001.
    nsq                    Runs nsq dockers.
    deps                   Install all dependencies.
EOF
  exit 1
}

GO=${GO:-$(which go)}
export GOPATH="$PWD"
export PATH="$PATH:$GOPATH/bin"


install_deps() {
  echo "Installing dependencies..."

  export GOPATH="$PWD/web"
  $GO get github.com/bitly/go-nsq
  $GO get github.com/mateuszdyminski/glog

  export GOPATH="$PWD/backend"
  $GO get github.com/bitly/go-nsq
  $GO get github.com/gorilla/websocket

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

startNsq() {
  echo "Starting nsq dockers..."
  set -e
  docker run -p 4150:4150 -p 4151:4151 -d quay.io/mateuszdyminski/nsq nsqd --broadcast-address=172.17.42.1 --lookupd-tcp-address=172.17.42.1:4160
  docker run -p 4160:4160 -p 4161:4161 -d quay.io/mateuszdyminski/nsq nsqlookupd --broadcast-address=172.17.42.1
  set +e
  echo "Nsq started with success!"
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
    exec "$PWD/backend/bin/backend" --p="8090" --nsqLookupd="127.0.0.1:4161"
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
  nsq)
    startNsq
  ;;
  *)
    usage
  ;;
esac
