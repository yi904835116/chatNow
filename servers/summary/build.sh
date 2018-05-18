set -e

export SUMMARY_CONTAINER=info-344-summary

GOOS=linux go build

docker build -t yi904835116/$SUMMARY_CONTAINER .

if [ "$(docker ps -aq --filter name=$SUMMARY_CONTAINER)" ]; then
    docker rm -f $SUMMARY_CONTAINER
fi

if [ "$(docker images -q -f dangling=true)" ]; then
    docker rmi $(docker images -q -f dangling=true)
fi

go clean