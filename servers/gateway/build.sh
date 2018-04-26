
set -e

export GATEWAY_CONTAINER=info344-gateway

GOOS=linux go build

docker build -t yi904835116/info344-server .




if [ "$(docker ps -aq --filter name=$GATEWAY_CONTAINER)" ]; then
    docker rm -f $GATEWAY_CONTAINER
fi


go clean