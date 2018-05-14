set -e

./build.sh


export CONTAINER_NAME=info344-mysqlDB

export MYSQL_DATABASE=info_344

export MYSQL_ROOT_PASSWORD="ABCD1234"


echo "mysql root password:" $MYSQL_ROOT_PASSWORD


