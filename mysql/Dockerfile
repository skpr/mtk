FROM debian:bullseye-slim

ENV DEBIAN_FRONTEND=noninteractive

RUN groupadd -r mysql && \
    useradd -r -g mysql mysql

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
      gnupg && \
    rm -rf /var/lib/apt/lists/*

RUN apt-key adv --keyserver hkp://keyserver.ubuntu.com --recv-keys 859BE8D7C586F538430B19C2467B942D3A79BD29 && \
    echo "deb http://repo.mysql.com/apt/debian/ bullseye mysql-8.0" | tee -a /etc/apt/sources.list.d/mysql.list

RUN { \
  echo mysql-community-server mysql-community-server/data-dir select ''; \
  echo mysql-community-server mysql-community-server/root-pass password 'root'; \
  echo mysql-community-server mysql-community-server/re-root-pass password 'root'; \
  echo mysql-community-server mysql-community-server/remove-test-db select false; \
} | debconf-set-selections && \
    apt-get update && \
    apt-get install -y \
      mysql-community-client \
      mysql-community-server && \
    rm -rf /var/lib/apt/lists/*

COPY scripts /usr/local/bin
RUN chmod +x /usr/local/bin/*

ENV MYSQL_TCP_PORT=3306

ADD my.cnf /etc/mysql/mysql.conf.d/mysqld.cnf

USER mysql

CMD ["/usr/local/bin/mysql-start"]
