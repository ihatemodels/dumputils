## Welcome to Dump Utils

**`IT is NOT 🚀 blazing-fast but still works!`**

# Under development and NOT ready to use!

## Table of Contents

**1. [Introduction](#Introduction)**  
**2. [Theory](#Theory)**  
**3. [Features](#Features)**  
**4. [Configuration Examples](<#configuration-examples>)**

- **[Simple PostgresSQL](<#single-postgres-database>)**
- **[Full Configuration Example](<#full-configuration-example>)**
- **[Examples Folder](docs/EXAMPLES.md)**

**5. [Orchestration using Kubernetes](#todo)**  
**6. [Manual Install](docs/INSTALL.md)**

### Introduction

This projects aims to pack and abstract the usage of all common dump utils such as all recent versions
of `pg_dump & pg_dumpall` , `mysqldump`, `mongodump` and more in a single docker container. The goal is to provide
orchestration (using Kubernetes) and observability of different backup procedures while still simple to use and
configure.

**<p style='color:orange'>WARN:</p> This tool can NOT be used as a standalone software without the proper binaries
like `pg_dump`, `mysqldump`, etc... The docker containers used in the examples contains all you need. If you want to use
this tool in Virtual Machines or other Environments you can follow the installation steps explained
in [INSTALL.md](docs/INSTALL.md)**

Full detailed examples of all kind of backups, notifications and supported storage outputs can be found in
the [Example Folder](example).

### Theory

Replication handles many error cases but by far not all. What about an admin that accidentally deleted the whole Prod
Cluster because they thought they were on dev? What about ransomware that affected all of your online systems? Having
cold snapshot of your data **desired on other place than your data itself** is the way to remain calm. Using all of
these dump tools is easy till the point where you have to scale with many databases and clusters. Then as with any
modern problem you need some kind of orchestration and observability. This is why I started this project.

## Features

1. **Clean Documentation**
    - [x] [How To Install Outside Docker](docs/INSTALL.md)
    - [x] [Restore Procedures](docs/RESTORE.md)
2. **Base Container to use manually for dump/restore procedures**
3. **Orchestration Strategies Using Kubernetes**
4. **Supported Dump Targets**
    - [x] PostgresSQL version 10,11,12,13,14
        - [x] Dump Single Database
        - [x] Dump All Databases In Different dump files
        - [x] Dump The Whole Server with `pg_dumpall`
    - [ ] MariaDB
    - [ ] MongoDB
    - [ ] Kafka Topics
    - [ ] Cassandra
    - [ ] Redis
5. **Supported Outputs**
    - [ ] Local Filesystem
    - [ ] S3 Compatible Storage System
    - [ ] SFTP Servers
    - [ ] RSYNC
6. **Supported Notification Channels**
    - [ ] Email
    - [ ] Slack
    - [ ] Matrix/Element
    - [ ] Telegram
7. **Support switch on/off concurrent backups based on configuration**

## Configuration Examples

It is important to know that each section of the configuration file can be disabled and enabled with the `enabled` flag.
If we add a new section, we must add the enabled: true field. You can see examples below. A WARN log message will be
printed if there are sections defined in the config, but `enabled` is not true. This was done because it is better to
change one flag than to comment on an entire section. In the context of Kubernetes configs-maps, this allows to manage
different configurations in more elegant way.

### Single Postgres Database

In this example we will backup the `sales` database running on 192.168.1.5:5432 ( Postgres Database Server version 11 ).
Saving the output to our file system in the `/opt/backup/` directory. We will receive Slack notification when the backup
is done or if it fails. We will use docker, so we don't have to worry about pg_tools installations as everything is
packed in the docker container.

- We will create `config.yaml` in our current directory with this content

```yaml
log:
  # choices: json|human
  type: human
  # choices: debug|info|warn|error
  level: "info"

databases:
  # dump single database example
  - name: "sales-east"
    enabled: true
    host: 192.168.1.5
    port: 5432
    database: sales
    username: postgres
    password: postgres
    version: 11
    # enabling this will pass the -v flag to pg_dump
    # all collected output will be printed as info log message
    verbose: true

output:
  filesystem:
    enabled: true
    # the user which starts dumputils must have write access
    # to this destination
    path: /opt/backups/

notifiers:
  slack: # TODO: add example
```

- Run the container

```bash
user@host: docker run --rm --name dumputils -h dumputils \ 
            -v /opt/backups:/opt/backups \ 
            -v $(pwd)/config.yaml:/etc/dumputils.config.yaml \ 
            -e DUMPUTILS_CONFIG_PATH="/etc/dumputils.config.yaml" \
            -it ihatemodels1/dump-utils:latest 

time=2022-06-03T01:49:36+03:00 level=info message=dumputils started
...
```

### Full configuration example