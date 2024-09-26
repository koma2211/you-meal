if you want to configure app for docker, change config path into invironment variable to configs/main.yaml.

# docker run -d \
#   --name postgres \
#   -e POSTGRES_PASSWORD=secret \
#   -e POSTGRES_DB=you_mealdb \
#   -p 5439:5432 \
#   postgres:16.4-alpine3.20
