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

    if [ -z $DATABASE_HOST ]; then
        fatal "Not found: DATABASE_HOST"
    fi

    if [ -z $DATABASE_PORT ]; then
        fatal "Not found: DATABASE_PORT"
    fi

    if [ -z $DATABASE_USER ]; then
        fatal "Not found: DATABASE_USER"
    fi

    if [ -z $DATABASE_PASSWORD ]; then
        fatal "Not found: DATABASE_PASSWORD"
    fi

    if [ -z $DATABASE_NAME ]; then
        fatal "Not found: DATABASE_NAME"
    fi

    if [ -z $FILE ]; then
        fatal "Not found: FILE"
    fi

    # Connection string which will be used to perform operations on the MySQL database.
    CONNECTION_STRING="mysql --host=$DATABASE_HOST --port=$DATABASE_PORT --user=$DATABASE_USER --password=$DATABASE_PASSWORD $DATABASE_NAME"

    TABLES=$($CONNECTION_STRING -e 'show tables' | awk '{ print $1}' | grep -v '^Tables' )

    if [ "$TABLES" == "" ]
    then
        fatal "Error - No tables found in $DATABASE_NAME database!"
    fi
    
    for table in $TABLES
    do
        info "Deleting $DATABASE_NAME/$table"
        $CONNECTION_STRING -e "drop table $table"
    done

    info "Restoring Database: $FILE"
    $CONNECTION_STRING < $FILE

    info "Restore Complete"
fi