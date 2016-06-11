#++++++++++++++++++++++++++++++++++++++
# Symbol-Exc Docker container in Alpine
#++++++++++++++++++++++++++++++++++++++

FROM golang:1.6-alpine
LABEL vendor=Katuyo
MAINTAINER palmtale<m@glad.so>

ENV SYMBOL_HOME=$GOPATH/src/github.com/katuyo/symbol-exchange

RUN set -ex \
    && apk add --no-cache git \
    && go get github.com/katuyo/symbol-exchange \
    && cd $SYMBOL_HOME \
    && go build \
    && apk del git && rm -rf /var/cache/apk/*

WORKDIR $SYMBOL_HOME
EXPOSE 4000
ENTRYPOINT ["./symbol-exchange"]
