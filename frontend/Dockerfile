FROM ubuntu:trusty

MAINTAINER mateuszdyminski

ADD app /home/app
ADD bin/server /usr/bin/server

RUN chmod 766 /home/app -R

CMD /usr/bin/server --dir=/home/app --p=9001
