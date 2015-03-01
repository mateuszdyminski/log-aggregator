FROM ubuntu:trusty

RUN apt-get update && apt-get install nginx curl wget -y && wget https://github.com/kelseyhightower/confd/releases/download/v0.7.1/confd-0.7.1-linux-amd64 && mv confd-0.7.1-linux-amd64 /usr/local/bin/confd && chmod +x /usr/local/bin/confd

ADD confd /etc/confd
ADD bin/confd-watch /usr/local/bin/confd-watch
COPY conf/nginx.conf /etc/nginx/nginx.conf

RUN chmod +x /usr/local/bin/confd-watch
RUN rm /etc/nginx/sites-enabled/default

WORKDIR /data

CMD ["bash"]
