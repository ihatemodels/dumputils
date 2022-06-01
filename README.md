# pgtools

## Using the base container as backup machine 

- Run the container

```bash
docker run --rm --name pgtools -h pgtools -it ihatemodels1/pgtools-base:latest /bin/bash
root@pgtools: ....
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
