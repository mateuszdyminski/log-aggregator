# Pull base image.
FROM ubuntu:trusty

# Install NSQ.
RUN \
  mkdir -p /tmp/nsq && apt-get update && apt-get install wget -y && \
  wget https://s3.amazonaws.com/bitly-downloads/nsq/nsq-0.3.2.linux-amd64.go1.4.1.tar.gz -O - | tar -xvz --strip=1 -C /tmp/nsq && \
  mv /tmp/nsq/bin/* /usr/local/bin/ && \
  rm -rf /tmp/nsq

# configure confd
ADD confd/bin/confd-0.7.1-linux-amd64 /usr/local/bin/confd
ADD confd /etc/confd
ADD run/confd-watch-nsqd /usr/local/bin/confd-watch-nsqd

RUN mkdir /etc/nsq
RUN chmod +x /usr/local/bin/confd

# Define working directory.
WORKDIR /data

# Define default command.
CMD ["bash"]

# Expose ports.
EXPOSE 4150
EXPOSE 4151
EXPOSE 4160
EXPOSE 4161
EXPOSE 4171
