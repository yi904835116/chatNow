docker build -t yi904835116/info344-client .


if [ "$(docker ps -aq --filter name=info344-client)" ]; then
    docker rm -f info-344-client
fi
