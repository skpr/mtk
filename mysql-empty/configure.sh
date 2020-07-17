#!/bin/bash

# A script for preconfiguring MySQL for local development
# which does not require a database.

USER=$1
PASS=$2
NAME=$3

mysql-start &
while ! nc -z 127.0.0.1 3306; do
  sleep 0.1
done

user-create $NAME $USER $PASS

mysql-stop
