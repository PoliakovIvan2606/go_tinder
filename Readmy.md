```sh
psql -U postgres -d postgres -h localhost -p 5432 -f init_db.sql
```

```sh
docker run --name postgres \
 -e POSTGRES_USER=postgres \
 -e POSTGRES_PASSWORD=postgres \
 -e POSTGRES_DB=postgres \
 -p 5432:5432 \
 -d postgres:latest
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
