set -e

./build.sh


export CONTAINER_NAME=info344-mysqlDB

export MYSQL_DATABASE=info_344

export MYSQL_ROOT_PASSWORD="ABCD1234"

export MYSQL_IMAGE=yi904835116/info344-mysql

echo "mysql root password:" $MYSQL_ROOT_PASSWORD



docker run -d \
--name $CONTAINER_NAME \
-p 3306:3306 \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_DATABASE=$MYSQL_DATABASE \
mysql