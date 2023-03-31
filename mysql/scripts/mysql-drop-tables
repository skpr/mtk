#!/bin/bash
set -euo pipefail

#/ Usage:       mysql-drop-tables [username] [password] [hostname] [database]
#/ Description: Drops tables for a database.
#/ Options:
#/   --help: Display this help message
usage() { grep '^#/' "$0" | cut -c4- ; exit 0 ; }
expr "$*" : ".*--help" > /dev/null && usage

info()    { echo "[INFO]  $*" ; }
fatal()   { echo "[FATAL] $*" ; exit 1 ; }

if [[ "${BASH_SOURCE[0]}" = "$0" ]]; then
    DATABASE_USERNAME="$1"
    DATABASE_PASSWORD="$2"
    DATABASE_HOSTNAME="$3"
    DATABASE_NAME="$4"

    if [ -z $DATABASE_HOSTNAME ]; then
        fatal "Not found: hostname"
    fi

    if [ -z $DATABASE_USERNAME ]; then
        fatal "Not found: username"
    fi

    if [ -z $DATABASE_PASSWORD ]; then
        fatal "Not found: password"
    fi

    if [ -z $DATABASE_NAME ]; then
        fatal "Not found: name"
    fi

    # Connection string which will be used to perform operations on the MySQL database.
    CONNECTION_STRING="mysql --host=$DATABASE_HOSTNAME --user=$DATABASE_USERNAME --password=$DATABASE_PASSWORD $DATABASE_NAME"

    # Grep is wrapped in { COMMAND || true; } because it returns an exit code of 1 if no results are found.
    # We have set -e for this script so if we don't do this wrapping the script will exit right here if no tables are found.
    # This script will fail on new environments with an empty database and break initial imports.
    TABLES=$($CONNECTION_STRING -e 'show full tables where Table_Type = "BASE TABLE"' | awk '{ print $1}' | { grep -v '^Tables' || true; } )

    if [ "$TABLES" != "" ]
    then
        for table in $TABLES
        do
            info "Deleting $DATABASE_NAME/$table"
            $CONNECTION_STRING -e "drop table $table"
        done
    fi

    # Grep is wrapped in { COMMAND || true; } because it returns an exit code of 1 if no results are found.
    # We have set -e for this script so if we don't do this wrapping the script will exit right here if no tables are found.
    # This script will fail on new environments with an empty database and break initial imports.
    VIEWS=$($CONNECTION_STRING -e 'show full tables where Table_Type = "VIEW"' | awk '{ print $1}' | { grep -v '^Tables' || true; } )

    if [ "$VIEWS" != "" ]
    then
        for view in $VIEWS
        do
            info "Deleting view $DATABASE_NAME/$view"
            $CONNECTION_STRING -e "drop view $view"
        done
    fi
fi
