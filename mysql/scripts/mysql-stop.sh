#!/bin/bash

PID=/run/mysqld/mysqld.pid

echo "Stopping Mysql..."

kill $(cat ${PID})