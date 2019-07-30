#!/bin/bash

PID=/run/mysqld/mysqld.pid

echo "Starting Mysql..."

mysqld --explicit_defaults_for_timestamp \
       --user=root \
       --pid-file=${PID}