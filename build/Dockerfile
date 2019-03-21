FROM golang:1.10.4-alpine

ARG environment=testoce
ARG zookeeperlist

MAINTAINER weiwenwang

ENV ENVIRONMENT=$environment
ENV PRJ_START_DELAY=5

COPY ./bin magneto/bin
COPY ./api magneto/api
COPY ./config magneto/config
COPY ./web magneto/web

EXPOSE 9308

ENTRYPOINT /bin/sh -c "echo 10.1.10.35 apollo-dev.kb.com >> /etc/hosts" \
    		&& /bin/sh -c "echo 10.1.10.35 apollo-fat.kb.com >> /etc/hosts" \
    		&& cd magneto \
    		&& mkdir -p status/magneto \
    		&& mkdir -p log \
    		&& sed -i -E "s/brokerlist.*=.*/brokerlist=\"${PRJ_BROKER_LIST}\"/g" ./config/magneto.${ENVIRONMENT}.toml \
    		&& zookeeperlist=$(echo "${PRJ_ZOOKEEPER_LIST}" | awk '{split($0,ips,",")} {for(i in ips){ print "\""ips[i]"\"," }}') \
    		&& sed -i -E "s/zookeeperAddresses.*=.*\[.*\]/zookeeperAddresses=[$zookeeperlist]/g" ./config/magneto.${ENVIRONMENT}.toml \
    		&& redislist=$(echo "${PRJ_REDIS_LIST}" | awk '{split($0,ips,",")} {for(i in ips){ print "\""ips[i]"\"," }}') \
            && sed -i -E "s/redis_svr.*=.*\[.*\]/redis_svr=[$redislist]/g" ./config/magneto.${ENVIRONMENT}.toml \
    		&& sleep ${PRJ_START_DELAY}s \
    		&& ./bin/magneto -config ./config/magneto.${ENVIRONMENT}.toml