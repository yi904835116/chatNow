
./build.sh

docker push yi904835116/info344-server

# Send run.sh to the cloud running remotely.
ssh -oStrictHostKeyChecking=no root@138.68.42.198 'bash -s' < run.sh
