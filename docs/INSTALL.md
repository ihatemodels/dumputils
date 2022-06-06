## Install 

This document refers to the installation steps of all related tools and the project itself if you want to use it outside docker, or you want to build your own containers.

### Ubuntu / Debian based distros

- Postgres Tools

```bash
# Install Base Utils
sudo apt-get update -y
sudo apt-get install -y --no-install-recommends curl gpg gnupg2 software-properties-common apt-transport-https ca-certificates
# Add Postgres Repositories
curl -fsSL https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo gpg --dearmor -o /etc/apt/trusted.gpg.d/postgresql.gpg

# Install Postgres Binaries
sudo apt-get update
sudo apt-get install -y --no-install-recommends \
      postgresql-client-10 postgresql-client-11 \
      postgresql-client-12 postgresql-client-13 postgresql-client-14
      
# Confirm binaries dirs
➜ tree /usr/lib/postgresql/
/usr/lib/postgresql/
├── 10
│   ├── bin
│   │   ├── clusterdb
│   │   ├── createdb
│   │   ├── createuser
│   │   ├── dropdb
│   │   ├── dropuser
│   │   ├── pg_basebackup
│   │   ├── pg_config
│   │   ├── pg_dump
│   │   ├── pg_dumpall
│   │   ├── pg_isready
│   │   ├── pg_receivewal
│   │   ├── pg_recvlogical
│   │   ├── pg_restore
│   │   ├── psql
│   │   ├── reindexdb
│   │   └── vacuumdb
│   └── lib
│       └── pgxs
│           ├── config
│           │   ├── install-sh
│           │   └── missing
│           └── src
│               ├── Makefile.global
│               ├── Makefile.port
│               ├── Makefile.shlib
│               ├── makefiles
│               │   └── pgxs.mk
│               ├── nls-global.mk
│               └── test
│                   ├── perl
│                   │   ├── PostgresNode.pm
│                   │   ├── RecursiveCopy.pm
│                   │   ├── SimpleTee.pm
│                   │   └── TestLib.pm
│                   └── regress
│                       └── pg_regress
├── 11
│   ├── bin
│   │   ├── clusterdb
│   │   ├── createdb
│   │   ├── createuser
│   │   ├── dropdb
....
```

