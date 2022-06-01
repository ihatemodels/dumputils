FROM golang:1.18.2-bullseye
ENV DEBIAN_FRONTEND noninteractive
LABEL maintainer="ihatemodels1@proton.me"

# Install Base Packages
RUN apt-get update -y && \
    apt-get install -y --no-install-recommends \
    curl gpg gnupg2 software-properties-common apt-transport-https ca-certificates

# Add PostgreSQL Repositories
RUN cd /tmp/ && \
    curl -fsSL https://www.postgresql.org/media/keys/ACCC4CF8.asc | gpg --dearmor -o /etc/apt/trusted.gpg.d/postgresql.gpg && \
    echo "deb http://apt.postgresql.org/pub/repos/apt/ bullseye-pgdg main" | tee /etc/apt/sources.list.d/pgdg.list && \
    apt-get update -y && \
    apt-get install -y --no-install-recommends \
    postgresql-client-10 postgresql-client-11 \
    postgresql-client-12 postgresql-client-13 postgresql-client-14 && \
    # install minio
    wget https://dl.min.io/client/mc/release/linux-amd64/mc && \
    mv mc /usr/local/bin/mc && chmod +x /usr/local/bin/mc && \
    rm -rf /tmp/*