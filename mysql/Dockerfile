FROM ubuntu:22.04

ENV TZ=Etc/UTC
ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
    && apt-get -y install \
      tzdata \
      bash \
      mysql-client \
      mysql-server \
    && apt-get clean \
    && rm -Rf /var/lib/apt/lists/* \
    && rm -Rf /usr/share/doc && rm -Rf /usr/share/man

RUN mkdir -p /opt/mysql/mysql
ADD my.cnf /opt/mysql/mysql/etc/my.cnf
RUN mysqld \
      --initialize \
      --user=mysql  \
      --basedir=/opt/mysql/mysql  \
      --datadir=/opt/mysql/mysql/data

RUN mkdir -p /run/mysqld && \
    chown mysql:mysql /run/mysqld
VOLUME /run/mysqld

COPY scripts /usr/local/bin
RUN chmod +x /usr/local/bin/*

ENV MYSQL_TCP_PORT=3306

USER mysql

CMD ["/usr/local/bin/mysql-start"]

RUN mysqld --validate-config