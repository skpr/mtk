#!/bin/bash

USER=$1
PASS=$2
NAME=$3
FILE=$4

mysql-start &
while ! nc -z localhost 3306; do
  sleep 0.1
done

database-create $NAME $FILE

user-create $NAME $USER $PASS

mysql-stop