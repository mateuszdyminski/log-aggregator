FROM ubuntu:trusty

MAINTAINER mateuszdyminski

ADD bin/backend /usr/bin/backend

CMD /usr/bin/backend --p=$PORT --nsqLookupd=$NSQ_LOOKUPD
