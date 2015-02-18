FROM dockerfile/nsq

RUN wget https://github.com/kelseyhightower/confd/releases/download/v0.6.3/confd-0.6.3-linux-amd64 && \
    mv confd-0.6.3-linux-amd64 /usr/local/bin/confd && \
    chmod +x /usr/local/bin/confd

RUN mkdir /etc/nsq

ADD confd /etc/confd
ADD bin/confd-watch-nsqd /usr/local/bin/confd-watch-nsqd

WORKDIR /data

CMD ["bash"]
