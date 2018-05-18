
set -e

export MESSAGING_CONTAINER=info-344-messaging

docker build -t yi904835116/$MESSAGING_CONTAINER .

if [ "$(docker ps -aq --filter name=$MESSAGING_CONTAINER)" ]; then
    docker rm -f $MESSAGING_CONTAINER
fi

# Remove dangling images.
if [ "$(docker images -q -f dangling=true)" ]; then
    docker rmi $(docker images -q -f dangling=true)
fi