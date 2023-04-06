
```

docker exec -it postgresql psql -U authuser
docker exec -i postgresql psql -U authuser -d auth < init.sql
```


```
DATABASE_URL: postgresql://authuser:mysecretpassword@localhost/auth?sslmode=disable
JWT_SECRET: secret
PORT: 1323
```


```
docker build -t microservice-auth .
docker run -d -p 1323:1323 --name microservice-auth -e DATABASE_URL='postgresql://authuser:mysecretpassword@172.17.0.2/auth?sslmode=disable' -e JWT_SECRET=secret -e PORT=1323 microservice-auth
```

```
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```