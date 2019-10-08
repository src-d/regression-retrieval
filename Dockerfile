FROM golang:1.12-stretch
MAINTAINER source{d}

ENV LOG_LEVEL=debug
ENV REG_REPOS=/cache/repos
ENV REG_BINARIES=/cache/binaries
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org

RUN apt-get update && \
    apt-get install -y dumb-init libonig-dev \
      git make bash gcc libxml2-dev && \
    apt-get autoremove -y && \
    ln -s /usr/local/go/bin/go /usr/bin

ADD build/regression-retrieval_linux_amd64/regression-retrieval /bin/

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/bin/regression-retrieval", "latest", "remote:master"]
