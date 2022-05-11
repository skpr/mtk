FROM alpine:3.15

RUN apk add bash \
            mariadb \
            mariadb-client

ADD my.cnf /etc/my.cnf.d/mariadb-server.cnf

RUN mysql_install_db --user=mysql --datadir=/var/lib/mysql

RUN mkdir -p /run/mysqld && \
    chown mysql:mysql /run/mysqld
VOLUME /run/mysqld

COPY scripts /usr/local/bin
RUN chmod +x /usr/local/bin/*

ENV MYSQL_TCP_PORT=3306

USER mysql

CMD ["/usr/local/bin/mysql-start"]
