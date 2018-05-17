set -e

./build.sh

docker push yi904835116/info-344-summary

export SERVER_IP=138.68.42.198

ssh -oStrictHostKeyChecking=no root@$SERVER_IP 'bash -s' < run.sh