## Welcome to Dump Utils 

**`IT is NOT 🚀 blazing-fast but still works!`**

This projects aims to pack and abstract the usage of all common dump utils such as all recent versions of `pg_dump & pg_dumpall` , `mysqldump`, `mongodump` and more in a single docker container. The goal is to provide orchestration (using Kubernetes) and observability of different backup procedures while still simple to use and configure.   

Full detailed examples of all kind of backups, notifications and supported storage output this tool provides can be found in the [Example Folder](example).

### Starting Simple

- ##### Single Postgres Database
In this example we will backup the `sales` database running on 192.168.1.5:5432 ( Postgres Database Server version 11 ). Saving the output to our file system in the `/opt/backup/` directory. We will receive Slack notification when the backup is done or if it fails. We will use docker, so we don't have to worry about pg_tools installations as everything is packed in the docker container.  

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

### Using the base container as backup machine 

- Run the container

```bash
docker run --rm --name pgtools -h pgtools -it ihatemodels1/pgtools-base:latest /bin/bash
```

- Create backup

```bash
root@pgtools: /usr/lib/postgresql/14/bin/pg_dump -h 192.168.1.5 -p 5432 -U postgres -Fc -Z 9 sales -f sales.dump
```

- Upload it to `minio`

```bash
root@pgtools: mc alias set minio-onprem http://192.168.1.51 {access_key} {secret_key}
# optionally create bucker
root@pgtools: mc mb minio-onprem/mybucket
root@pgtools: mc cp sales.dump minio-onprem/mybucket
```

- Restore it on other server

```bash
root@pgtools: /usr/lib/postgresql/14/bin/pg_restore -h 192.168.1.6 -p 5432 -U postgres -d sales_old sales.dump
```
