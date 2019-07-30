#!/bin/bash
set -euo pipefail

#/ Usage:       database-backup.sh > file.sql
#/ Description: Dumps MySQL database as part of a backup process.
#/ Options:
#/   --help: Display this help message
usage() { grep '^#/' "$0" | cut -c4- ; exit 0 ; }
expr "$*" : ".*--help" > /dev/null && usage

info()    { echo "[INFO]  $*" ; }
fatal()   { echo "[FATAL] $*" ; exit 1 ; }

if [[ "${BASH_SOURCE[0]}" = "$0" ]]; then
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

    info "Backup Started"

    mysqldump --single-transaction \
            --host=$DATABASE_HOST \
            --port=$DATABASE_PORT \
            --user=$DATABASE_USER \
            --password=$DATABASE_PASSWORD $DATABASE_NAME

    info "Backup Complete"
fi