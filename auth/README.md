
```

docker exec -it postgresql psql -U authuser
docker exec -i postgresql psql -U authuser -d auth < init.sql
```