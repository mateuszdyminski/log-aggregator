FROM ubuntu:trusty

MAINTAINER mateuszdyminski

ADD bin/main /usr/bin/main

CMD main --logtostderr=1 --toChannel=true --nsqHost=nsqWorker --p=$PORT
