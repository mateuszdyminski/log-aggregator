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
    build                   Builds all.
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
  docker build -t quay.io/mateuszdyminski/frontend:"$VERSION" "$PWD/frontend/."
  containerId=$(docker run -d quay.io/mateuszdyminski/frontend:"$VERSION")
  docker commit "$containerId" quay.io/mateuszdyminski/frontend
  docker push quay.io/mateuszdyminski/frontend
  docker kill "$containerId"
  docker rm "$containerId"
  set +e
  echo "Frontend pushed to quay.io!"
}

pushBackend() {
  echo "Pushing backend container started..."
  buildFrontend
  set -e
  docker build -t quay.io/mateuszdyminski/backend:"$VERSION" "$PWD/backend/."
  containerId=$(docker run -d quay.io/mateuszdyminski/frontend:"$VERSION")
  docker commit "$containerId" quay.io/mateuszdyminski/backend
  docker push quay.io/mateuszdyminski/backend
  docker kill "$containerId"
  docker rm "$containerId"
  set +e
  echo "Backend pushed to quay.io!"
}

pushWeb() {
  echo "Pushing web container started..."
  buildFrontend
  set -e
  docker build -t quay.io/mateuszdyminski/web:"$VERSION" "$PWD/web/."
  containerId=$(docker run -d quay.io/mateuszdyminski/web:"$VERSION")
  docker commit "$containerId" quay.io/mateuszdyminski/web
  docker push quay.io/mateuszdyminski/web
  docker kill "$containerId"
  docker rm "$containerId"
  set +e
  echo "Web pushed to quay.io!"
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
  *)
    usage
  ;;
esac
