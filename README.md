start 
docker-compose -f deployments/docker-compose.yml up

order-service создаёт заказ → пишет в Kafka order.created

eta-service читает → считает ETA → публикует order.eta.updated

order-service подписывается на order.eta.updated и обновляет ETA в БД