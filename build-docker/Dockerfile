FROM scratch as client

ENV DOCKER_CLIENT_VERSION=19.03.8
ENV DOCKER_API_VERSION=1.40
ADD https://download.docker.com/linux/static/stable/x86_64/docker-${DOCKER_CLIENT_VERSION}.tgz .

FROM golang:1-alpine
ENV TZ=Asia/Tokyo
ARG ARG_DOCKER_CLIENT_VERSION=19.03.8

ENV DOCKER_CLIENT_VERSION=${ARG_DOCKER_CLIENT_VERSION}
ENV DOCKER_API_VERSION=1.40

# ENV PROTOBUF_VERSION 3.13.0
# ENV PROTOBUF_URL https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOBUF_VERSION/protobuf-cpp-$PROTOBUF_VERSION.tar.gz

# WORKDIR /tmp

RUN set -eux \
    && apk update \
    && apk --no-cache add tzdata gcc libc-dev git make bash curl build-base autoconf automake libtool
#     && curl -L -o /tmp/protobuf.tar.gz $PROTOBUF_URL \
#     && tar xzvf protobuf.tar.gz \
#     && cd /tmp/protobuf-$PROTOBUF_VERSION \
#     && ./autogen.sh \
#     && ./configure \
#     && make -j 3 \
#     && make install \
#     && go get -u github.com/golang/protobuf/protoc-gen-go \
#     && go get -u github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc \
#     && go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
#     && go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger


WORKDIR /tmp

COPY --from=client docker-${DOCKER_CLIENT_VERSION}.tgz .

RUN set -eux \
    && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
    && echo "Asia/Tokyo" >  /etc/timezone \
    && apk del tzdata \
    && tar xzvf docker-${DOCKER_CLIENT_VERSION}.tgz \
    && mv docker/* /usr/bin/ \
    && rm -rf docker-${DOCKER_CLIENT_VERSION}.tgz \
    && mkdir /src \
    && go get -u github.com/rubenv/sql-migrate/... \
    && rm  -rf /tmp/* /var/cache/apk/* \
    && git clone https://github.com/awslabs/git-secrets.git \
    && cd git-secrets \
    && make install \
    && git secrets --register-aws --global \
    && cd ../ \
    && rm -rf git-secrets

WORKDIR /src/

CMD [ "sh" ]