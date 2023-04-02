
```
docker build -t microservice-gateway .
docker run -d -p 1324:1324 --name microservice-gateway -e MONGO_DB_URL=mongodb://root:password@172.17.0.3:27017 -e AUTH_SERVICE_URL=http://172.17.0.4:1323 -e PORT=1324 -e RABBITMQ_URL=amqp://guest:guest@172.17.0.5:5672 microservice-gateway
```