FROM java:alpine
MAINTAINER ShiYi <me@shiyicode.com>

RUN apk upgrade --update;
RUN apk add linux-headers;
RUN apk add bash;
RUN apk add git;
RUN apk add curl;
RUN apk add python;
RUN apk add g++;
RUN apk add gcc;
RUN apk add libc-dev;
RUN curl https://storage.googleapis.com/golang/go1.9.linux-amd64.tar.gz | tar xzf - -C /; \
    mv /go /goroot;

RUN mkdir /lib64; \
    ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2;

ENV GOPATH /go
ENV GOROOT /goroot
ENV GOBIN=$GOPATH/bin
ENV PATH=$PATH:$GOROOT/bin

RUN mkdir -p /go/src; \
    cd /go/src; \
    git clone https://github.com/open-fightcoder/oj-judger.git; \
    cd oj-judger; \
    ./build.sh; \
    cd output; \
    ./control.sh start;

WORKDIR /go/src/oj-judger/output

CMD while true; do sleep 1; done