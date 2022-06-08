FROM golang:1.18.2-bullseye
ENV DEBIAN_FRONTEND noninteractive
LABEL maintainer="ihatemodels1@proton.me"

# Install Base Packages
RUN apt-get update -y && \
    apt-get install -y --no-install-recommends \
    curl wget gpg gnupg2 software-properties-common apt-transport-https ca-certificates

# Add PostgreSQL Repositories
RUN cd /tmp/ && \
    curl -fsSL https://www.postgresql.org/media/keys/ACCC4CF8.asc | gpg --dearmor -o /etc/apt/trusted.gpg.d/postgresql.gpg && \
    echo "deb http://apt.postgresql.org/pub/repos/apt/ bullseye-pgdg main" | tee /etc/apt/sources.list.d/pgdg.list && \
    curl -fsSl https://mariadb.org/mariadb_release_signing_key.asc | gpg --dearmor -o /etc/apt/trusted.gpg.d/mariadb.gpg && \
    echo "deb [arch=amd64,arm64,ppc64el] https://ftp.ubuntu-tw.org/mirror/mariadb/repo/10.6/debian bullseye main" | tee /etc/apt/sources.list.d/mariadb.list && \
    apt-get update -y && \
    apt-get install -y --no-install-recommends \
        postgresql-client-10 postgresql-client-11 \
        postgresql-client-12 postgresql-client-13 postgresql-client-14 \
        mariadb-client && \
    # install minio client
    wget https://dl.min.io/client/mc/release/linux-amd64/mc && \
    mv mc /usr/local/bin/mc && chmod +x /usr/local/bin/mc && \
    # install mongo tools
    wget https://fastdl.mongodb.org/tools/db/mongodb-database-tools-debian10-x86_64-100.5.2.deb -O mongo-tools.deb && \
    dpkg -i mongo-tools.deb && \
    rm -rf /tmp/*