set -e


export TLSCERT=/etc/letsencrypt/live/api.patrick-yi.com/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.patrick-yi.com/privkey.pem


# Make sure to get the latest image.
# pull most current version of example web site container image
docker pull yi904835116/info344-server

# stop and remove current container instance
docker rm -f info344-server


# Run Info 344 API Gateway Docker container inside our appnet private network.
docker run -d \
--name 344gateway \
-p 443:443 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
yi904835116/info344-server