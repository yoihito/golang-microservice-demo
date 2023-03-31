
```

docker exec -it postgresql psql -U authuser
docker exec -i postgresql psql -U authuser -d auth < init.sql
```


```
DATABASE_URL: postgresql://authuser:mysecretpassword@localhost/auth?sslmode=disable
JWT_SECRET: secret
PORT: 1323
```