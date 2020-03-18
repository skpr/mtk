#!/bin/bash

export MYSQL_HOSTNAME=127.0.0.1
export MYSQL_USERNAME=tmp
export MYSQL_PASSWORD=tmp
export MYSQL_DATABASE=tmp

SOURCE=$1
TARGET=$2

echo "Starting MySQL"
mysql-start &
while ! nc -z localhost 3306; do
  sleep 0.1
done

echo "Creating temporary database"
database-create $MYSQL_DATABASE $SOURCE

echo "Creating temporary account"
user-create $MYSQL_DATABASE $MYSQL_USERNAME $MYSQL_PASSWORD

echo "Dumping sanitized database"
mtk db dump > $TARGET

echo "Stopping MySQL"
mysql-stop