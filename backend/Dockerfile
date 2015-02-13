FROM ubuntu:trusty

MAINTAINER mateuszdyminski

ADD bin/backend /usr/bin/backend

CMD /usr/bin/backend --p=8080 --kafka=kafka --wait=10 --topic=$TOPIC
