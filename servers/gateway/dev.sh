export TLSCERT=/Users/zhaoyangyi/go/src/github.com/info344-s18/challenges-yi904835116/tls/fullchain.pem
export TLSKEY=/Users/zhaoyangyi/go/src/github.com/info344-s18/challenges-yi904835116/tlsprivkey.pem

docker run -d \                            #run as detached process
--name 344gateway \                        #name for container instance
-p 443:443 \                               #publish port 443
-v /etc/letsencrypt:/etc/letsencrypt:ro \  #mount /etc/letsencrypt as /etc/letsencrypt in the container, read-only
-e TLSCERT=$TLSCERT \                      #forward TLSCERT env var into container
-e TLSKEY=$TLSKEY \                        #forward TLSKEY env var into container
your-dockerhub-name/your-container-name    #name of container image