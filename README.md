if you want to configure app for docker, change config path into invironment variable to configs/main.yaml.

docker run -d \
   --name postgres \
   -e POSTGRES_USER=root \
   -e POSTGRES_PASSWORD=secret \
   -e POSTGRES_DB=you_mealdb \
   -p 5439:5432 \
   postgres:16.4-alpine3.20


<!-- docker run -d \
  --name redis \
  -p 6379:6379 \
  -e REDIS_PASSWORD=secret \
  redis:7.4.0-alpine \
  redis-server --requirepass secret -->