FROM ihatemodels1/dumputils-base:latest
WORKDIR /builder
RUN mkdir -p /tmp/
COPY . ./