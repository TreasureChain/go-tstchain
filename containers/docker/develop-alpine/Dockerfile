FROM alpine:3.7

RUN \
  apk add --update go git make gcc musl-dev linux-headers ca-certificates && \
  git clone --depth 1 https://github.com/TreasureChain/go-tstchain && \
  (cd go-tstchain && make gtst) && \
  cp go-tstchain/build/bin/gtst /gtst && \
  apk del go git make gcc musl-dev linux-headers && \
  rm -rf /go-tstchain && rm -rf /var/cache/apk/*

EXPOSE 8545
EXPOSE 30303

ENTRYPOINT ["/gtst"]
