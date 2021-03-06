# Copyright (c) 2018 nChain Ltd
# Distributed under the GNU GPL v3.0, see accompanying file LICENSE for details
# based on work by Adrian Macneil from https://github.com/amacneil/docker-bitcoin
FROM golang:stretch

ENV BCH_VERSION 0.17.1
ENV BCH_URL https://download.bitcoinabc.org/${BCH_VERSION}/linux/bitcoin-abc-${BCH_VERSION}-x86_64-linux-gnu.tar.gz
ENV BCH_SHA256 eccf8b61ba0549f6839e586c7dc6fc4bf6d7591ac432aaea8a7df0266b113d27

ADD $BCH_URL /tmp/bitcoin.tar.gz
RUN cd /tmp \
    && echo "$BCH_SHA256  bitcoin.tar.gz" | sha256sum -c - \
    && tar -xzvf bitcoin.tar.gz -C /usr/local --strip-components=1 --exclude=*-qt \
    && rm bitcoin.tar.gz

RUN addgroup bitcoin && adduser --gecos "" --home /home/bitcoin --disabled-password --ingroup bitcoin bitcoin
ENV BCH_DATA /data
# RUN mkdir "$BCH_DATA" \
#     && chown -R bitcoin:bitcoin "$BCH_DATA" \
#     && ln -sfn "$BCH_DATA" /home/bitcoin/.bitcoin \
#     && chown -h bitcoin:bitcoin /home/bitcoin/.bitcoin
VOLUME /data

RUN go get -u github.com/rohenaz/dtvcash
COPY config.yaml /home/bitcoin/go/src/github.com/rohenaz/dtvcash/config.yaml
RUN go build github.com/rohenaz/dtvcash
# dont need this bc go build it building it in the gopath?
# RUN mv /home/bitcoin/go/src/github.com/rohenaz/dtvcash/dtvcash /go/bin/
COPY entrypoint.sh /entrypoint.sh
# COPY bitcoin.conf /bitcoin.conf
USER bitcoin

EXPOSE 8332 8333 18332 18333
# CMD ["bitcoind"]

ENTRYPOINT ["/entrypoint.sh"]