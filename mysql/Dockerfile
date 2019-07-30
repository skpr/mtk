FROM alpine:latest

RUN apk add --update bash mysql mysql-client && \
    rm -f /var/cache/apk/* && \
    mkdir -p /run/mysqld
COPY my.cnf /etc/mysql/my.cnf

ENV MYSQL_TCP_PORT=3306

RUN mysql_install_db --user=root

# Scripts used for backup and restore tasks.
ADD scripts/database-backup.sh /usr/local/bin/database-backup
ADD scripts/database-restore.sh /usr/local/bin/database-restore

# Scripts used for database image building.
ADD scripts/database-create.sh /usr/local/bin/database-create
ADD scripts/database-import.sh /usr/local/bin/database-import
ADD scripts/mysql-start.sh /usr/local/bin/mysql-start
ADD scripts/mysql-stop.sh /usr/local/bin/mysql-stop
ADD scripts/user-create.sh /usr/local/bin/user-create

RUN chmod +x /usr/local/bin/*

CMD ["/usr/local/bin/mysql-start"]