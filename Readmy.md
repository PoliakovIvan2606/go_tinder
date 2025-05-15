```sh
psql -U postgres -d postgres -h localhost -p 5432 -f init_db.sql

docker exec -it my_redis redis-cli
```

```sh
docker run --name pg-postgis \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=postgres \
  -p 5432:5432 \
  -d postgis/postgis:16-3.4
```

```sh
docker run -d \
  --name my_redis \
  -p 6379:6379 \
  redis:7 \
  redis-server --appendonly yes
```

Создать файл .env и добавить в него такие переменные
```
AWS_ACCESS_KEY_ID=97788b4134484665a2cac42b8057ed71
AWS_SECRET_ACCESS_KEY=7b4237b02a3148309767453cab248501
```

Затем с помощью команды добавить их в окружение
```sh
export $(grep -v '^#' .env | xargs)
```

> Хотел предоставлять токены при создании пользователя, но куки не добавлялись при повторном запросе и пэтому /upload не проходил.