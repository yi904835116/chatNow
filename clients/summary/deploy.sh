./build.sh

docker push yi904835116/info344-client

ssh -oStrictHostKeyChecking=no root@159.65.74.235 'bash -s' < run.sh