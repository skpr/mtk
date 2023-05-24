FROM ubuntu:22.04

ENV TZ=Etc/UTC
ENV DEBIAN_FRONTEND=noninteractive

RUN groupadd -r mysql && useradd -r -g mysql mysql

RUN mkdir /docker-entrypoint-initdb.d

RUN apt-get update \
    && apt-get -y install --no-install-recommends  \
      bzip2 \
      tzdata \
      openssl \
      bash \
      mysql-client \
      mysql-server \
      perl \
      xz-utils \
      zstd \
    && apt-get clean \
    && rm -Rf /var/lib/apt/lists/* \
    && rm -Rf /usr/share/doc \
    && rm -Rf /usr/share/man \
    && rm -rf /var/lib/mysql \
    && mkdir -p /opt/mysql /var/lib/mysql /var/run/mysqld \
    && chown -R mysql:mysql /opt/mysql /var/lib/mysql /var/run/mysqld \
    && chmod 1777 /opt/mysql /var/lib/mysql /var/run/mysqld

VOLUME /var/lib/mysql

ADD my.cnf /etc/mysql/my.cnf
COPY scripts /usr/local/bin
RUN chmod +x /usr/local/bin/*
RUN mysqld --validate-config

USER mysql

RUN mysqld \
      --initialize \
      --user=mysql \
      --basedir=/opt/mysql  \
      --datadir=/opt/mysql/data

EXPOSE 3306 33060
CMD ["/usr/local/bin/mysql-start"]
