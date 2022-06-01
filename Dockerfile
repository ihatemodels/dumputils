FROM ihatemodels1/pgtools-base:latest
WORKDIR /builder
COPY . ./
RUN mkdir -p /tmp/ && \
    GOOS=linux go build -v -mod=vendor -o pgtools ./

WORKDIR /app/
RUN cp /builder/pgtools . && \
    rm -rf /builder

CMD ["./pgtools"]