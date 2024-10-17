```bash
docker run --name local_mysql -e MYSQL_ROOT_PASSWORD=admin111 -p 3306:3306 -d mysql:8.0

docker run --name local_redis -p 6379:6379 -d redis
```