#!/bin/bash

NAME=$1
USER=$2
PASS=$3

echo "Creating user..."

mysql -u root -e "CREATE USER '${USER}'@'%' IDENTIFIED BY '${PASS}';"
mysql -u root -e "GRANT ALL PRIVILEGES ON ${NAME}.* TO '${USER}'@'%' IDENTIFIED BY '${PASS}';"
