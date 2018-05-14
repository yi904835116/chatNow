set -e


export TLSCERT=/etc/letsencrypt/live/api.patrick-yi.com/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.patrick-yi.com/privkey.pem

export GATEWAY_CONTAINER=344gateway
export ADDR=:443

export MYSQL_CONTAINER=info344-mysqlDB
export DBNAME=info_344
export MYSQL_ROOT_PASSWORD="ABCD1234"

export REDIS_CONTAINER=redis-server
export REDISADDR=$REDIS_CONTAINER:6379

export APP_NETWORK=appnet

export SESSIONKEY=secretsigningkey

export MYSQL_DATABASE=info_344


export MYSQL_ADDR=$MYSQL_CONTAINER:3306
# export MYSQL_ADDR=127.0.0.1:3306

# dsn := fmt.Sprintf("root:%s@tcp(192.168.99.100:3306)/info_344", os.Getenv("MYSQL_ROOT_PASSWORD"))

# export DSN=root:$MYSQL_ROOT_PASSWORD@tcp(127.0.0.1:3306)/$DBNAME
# export DSN=root:$MYSQL_ROOT_PASSWORD@tcp$MYSQL_ADDR/$DBNAME




echo helloooooooooo

# Make sure to get the latest image.
# pull most current version of example web site container image
docker pull yi904835116/info344-server

docker pull yi904835116/info344-mysql


# Create Docker private network if not exist.
if ! [ "$(docker network ls | grep $APP_NETWORK)" ]; then
    docker network create appnet
fi

# # Remove the old containers first.
if [ "$(docker ps -aq --filter name=$GATEWAY_CONTAINER)" ]; then

    docker rm -f $GATEWAY_CONTAINER
fi

if [ "$(docker ps -aq --filter name=$REDIS_CONTAINER)" ]; then
    docker rm -f $REDIS_CONTAINER
fi


if [ "$(docker ps -aq --filter name=$MYSQL_CONTAINER)" ]; then
    docker rm -f $MYSQL_CONTAINER
fi



# Run MySQL Docker container
docker run -d \
--name $MYSQL_CONTAINER \
--network appnet \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_DATABASE=$MYSQL_DATABASE \
yi904835116/info344-mysql

# Run Redis Docker container inside our appnet private network.
docker run \
-d \
--name $REDIS_CONTAINER \
--network appnet \
--restart unless-stopped \
redis

docker run \
-d \
-p 443:443 \
--name $GATEWAY_CONTAINER \
--network appnet \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
-e SESSIONKEY=$SESSIONKEY \
-e ADDR=$ADDR \
-e REDISADDR=$REDISADDR \
-e MYSQL_ADDR=$MYSQL_ADDR \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_DATABASE=$MYSQL_DATABASE \
--restart unless-stopped \
yi904835116/info344-server
