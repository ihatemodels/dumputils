FROM ihatemodels1/pgtools-base:latest
WORKDIR /builder
RUN mkdir -p /tmp/
COPY . ./