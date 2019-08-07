#!/bin/bash
set -euo pipefail

#/ Usage:       database-backup.sh > file.sql
#/ Description: Dumps MySQL database as part of a backup process.
#/ Options:
#/   --help: Display this help message
usage() { grep '^#/' "$0" | cut -c4- ; exit 0 ; }
expr "$*" : ".*--help" > /dev/null && usage

echoerr() { printf "%s\n" "$*" >&2; }
info()    { echoerr "[INFO]  $*" ; }
fatal()   { echoerr "[FATAL] $*" ; exit 1 ; }

if [[ "${BASH_SOURCE[0]}" = "$0" ]]; then
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

    info "Backup Started"

    mysqldump --single-transaction \
            --host=$MYSQL_HOSTNAME \
            --port=$MYSQL_PORT \
            --user=$MYSQL_PASSWORD \
            --password=$MYSQL_PASSWORD $MYSQL_DATABASE

    info "Backup Complete"
fi
