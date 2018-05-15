

docker pull yi904835116/info344-client
docker rm -f 344client

export TLSCERT=/etc/letsencrypt/live/web.patrick-yi.com/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/web.patrick-yi.com/privkey.pem

docker run -d --name tmp-nginx nginx
docker cp tmp-nginx:/etc/nginx/conf.d/default.conf default.conf
docker rm -f tmp-nginx

docker run -d \
--name 344client \
-p 80:80 -p 443:443 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
yi904835116/info344-client