FROM ubuntu:trusty

MAINTAINER mateuszdyminski

ADD bin/backend /usr/bin/backend


# configure confd
ADD confd/bin/confd-0.7.1-linux-amd64 /usr/local/bin/confd
ADD confd /etc/confd
ADD run/confd-watch-nsqlookupd /usr/local/bin/confd-watch-nsqlookupd

RUN mkdir -p /etc/backend/ 
RUN chmod +x /usr/local/bin/confd

# Env
ENV PORT 8090

# Expose ports.
EXPOSE $PORT

# Define default command.
CMD ["bash"]