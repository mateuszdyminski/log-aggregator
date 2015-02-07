### RUN
export GOPATH=$PWD && go build -o bin/main emiter && bin/main --logtostderr=1 --toChannel=true
