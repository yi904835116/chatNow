./build.sh


export SERVER_IP=159.65.74.235

docker push yi904835116/info344-client



# ssh -oStrictHostKeyChecking=no root@159.65.74.235 'bash -s' < run.sh
ssh -oStrictHostKeyChecking=no root@$SERVER_IP 'bash -s' < run.sh