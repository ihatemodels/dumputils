## Restore 

This document refers to techniques that must be used to restore backup's created with `dumputils`

- **Postgres Single Database** 

In this example we will backup the `sales` database from server 192.168.1.3 and later restore it on database `sales_new` in server 192.168.1.6

- Used **dumputils** config

```yaml
databases:
  # dump single database example
  - name: "sales-east"
    host: 192.168.1.3
    port: 5432
    database: sales
    username: postgres
    password: postgres
    version: 14
    # enabling this will pass the -v flag to pg_dump
    # all collected output will be printed as info log message
    verbose: true
    dumpAll: false
    dumpServer: false
    
filesystem:
    enabled: true
    # the user which starts pgtools must have write access
    # to this destination
    path: /opt/backups/
```

- Restore procedure

`We asume that you will run this container on the same storage where you ran dumputils when you took the backup`

```bash
docker run --rm --name dumputils \
      -h dumputils \
      -v /opt/backups:/opt/backups
      -it ihatemodels1/dumputils-base:latest
      
root@dumputils: /usr/lib/postgresql/14/bin/createdb -h 192.168.1.6 -p 5432 -U postgres sales_new
Password: # Enter the postgres user password here
root@dumputils: cd /opt/backups
root@dumputils: /usr/lib/postgresql/14/bin/pg_restore -h 192.168.1.6 -p 5432 -U postgres -d sales_new sales-east-sales-2022-04-16-14-45-09.dump
```