set -e

export MESSAGING_CONTAINER=info-344-messaging
export MQ_CONTAINER=rabbitmq-server
export MONGO_CONTAINER=mongo-server
export APP_NETWORK=appnet
export DBNAME="info_344"

docker pull yi904835116/$MESSAGING_CONTAINER

if [ "$(docker ps -aq --filter name=$MESSAGING_CONTAINER)" ]; then
    docker rm -f $MESSAGING_CONTAINER
fi

if [ "$(docker ps -aq --filter name=$MONGO_CONTAINER)" ]; then
    docker rm -f $MONGO_CONTAINER
fi

if [ "$(docker images -q -f dangling=true)" ]; then
    docker rmi $(docker images -q -f dangling=true)
fi

if ! [ "$(docker network ls | grep $APP_NETWORK)" ]; then
    docker network create $APP_NETWORK
fi

# Run Mongo Docker container inside our appnet private network.
docker run \
-d \
--name mongo-server \
--network $APP_NETWORK \
--restart unless-stopped \
mongo


docker run \
-d \
-e ADDR=$MESSAGING_CONTAINER:80 \
-e MQADDR=$MQ_CONTAINER:5672 \
-e DBADDR=mongo-server:27017 \
-e REDISADDR=redis-server \
-e SUMMARYSVCADDR=info-344-summary:80 \
--name $MESSAGING_CONTAINER \
--network $APP_NETWORK \
--restart unless-stopped \
yi904835116/$MESSAGING_CONTAINER