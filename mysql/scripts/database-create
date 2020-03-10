#!/bin/bash

NAME=$1
FILE=$2

echo "Creating database..."
mysql -u root -e "CREATE DATABASE ${NAME}"

echo "Importing database..."
mysql -u root -D ${NAME} < ${FILE}