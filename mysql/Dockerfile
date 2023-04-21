FROM ubuntu:22.04

RUN apt-get update \
    && apt-get install -y \
      bash \
      mysql-client \
      mysql-server \
    && apt-get clean \
    && rm -Rf /var/lib/apt/lists/* \
    && rm -Rf /usr/share/doc && rm -Rf /usr/share/man

ADD my.cnf /etc/my.cnf

RUN mkdir -p /run/mysqld && \
    chown mysql:mysql /run/mysqld
VOLUME /run/mysqld

COPY scripts /usr/local/bin
RUN chmod +x /usr/local/bin/*

ENV MYSQL_TCP_PORT=3306

USER mysql

CMD ["/usr/local/bin/mysql-start"]

RUN mysqld --validate-config