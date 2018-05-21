set -e

export MESSAGING_CONTAINER=info-344-messaging

./build.sh

docker push yi904835116/$MESSAGING_CONTAINER

export SERVER_IP=138.68.42.198

ssh -oStrictHostKeyChecking=no root@$SERVER_IP 'bash -s' < run.sh