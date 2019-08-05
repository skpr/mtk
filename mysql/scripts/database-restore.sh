#!/bin/bash
set -euo pipefail

#/ Usage:       database-restore.sh file.sql
#/ Description: Drops tables and imports a database as part the restore process.
#/ Options:
#/   --help: Display this help message
usage() { grep '^#/' "$0" | cut -c4- ; exit 0 ; }
expr "$*" : ".*--help" > /dev/null && usage

info()    { echo "[INFO]  $*" ; }
fatal()   { echo "[FATAL] $*" ; exit 1 ; }

if [[ "${BASH_SOURCE[0]}" = "$0" ]]; then
    FILE=$1

    if [ -z $MYSQL_HOSTNAME ]; then
        fatal "Not found: MYSQL_HOSTNAME"
    fi

    if [ -z $MYSQL_PORT ]; then
        fatal "Not found: MYSQL_PORT"
    fi

    if [ -z $MYSQL_USERNAME ]; then
        fatal "Not found: MYSQL_USERNAME"
    fi

    if [ -z $MYSQL_PASSWORD ]; then
        fatal "Not found: MYSQL_PASSWORD"
    fi

    if [ -z $MYSQL_DATABASE ]; then
        fatal "Not found: MYSQL_DATABASE"
    fi

    if [ -z $FILE ]; then
        fatal "Not found: FILE"
    fi

    # Connection string which will be used to perform operations on the MySQL database.
    CONNECTION_STRING="mysql --host=$MYSQL_HOSTNAME --port=$MYSQL_PORT --user=$MYSQL_USERNAME --password=$MYSQL_PASSWORD $MYSQL_DATABASE"

    TABLES=$($CONNECTION_STRING -e 'show tables' | awk '{ print $1}' | grep -v '^Tables' )

    if [ "$TABLES" == "" ]
    then
        fatal "Error - No tables found in $MYSQL_DATABASE database!"
    fi
    
    for table in $TABLES
    do
        info "Deleting $MYSQL_DATABASE/$table"
        $CONNECTION_STRING -e "drop table $table"
    done

    info "Restoring Database: $FILE"
    $CONNECTION_STRING < $FILE

    info "Restore Complete"
fi