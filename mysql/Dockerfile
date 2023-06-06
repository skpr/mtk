FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive

RUN groupadd -r mysql && \
    useradd -r -g mysql mysql

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
      netcat \
      gnupg

RUN { \
  echo mysql-community-server mysql-community-server/data-dir select ''; \
  echo mysql-community-server mysql-community-server/root-pass password 'root'; \
  echo mysql-community-server mysql-community-server/re-root-pass password 'root'; \
  echo mysql-community-server mysql-community-server/remove-test-db select false; \
} | debconf-set-selections && \
    apt-get update && \
    apt-get install -y \
      mysql-client \
      mysql-server && \
    rm -rf /var/lib/apt/lists/*

RUN mkdir -p /run/mysqld && \
    chown mysql:mysql /run/mysqld
VOLUME /run/mysqld

COPY scripts /usr/local/bin
RUN chmod +x /usr/local/bin/*

ENV MYSQL_TCP_PORT=3306

ADD my.cnf /etc/mysql/mysql.cnf

USER mysql

CMD ["/usr/local/bin/mysql-start"]
