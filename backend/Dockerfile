FROM ubuntu:trusty

MAINTAINER mateuszdyminski

ADD bin/backend /usr/bin/backend

CMD /usr/bin/backend --p=8080 --nsqHost=nsqMaster --nsqTopic=all
