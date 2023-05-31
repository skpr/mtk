FROM skpr/mtk-mysql:v2-latest-amd64

ARG db_name=local
ARG db_user=local
ARG db_pass=local
ARG db_file=/tmp/db.sql

ARG X=y

ADD db.sql /tmp/db.sql

RUN database-import ${db_name} ${db_user} ${db_pass} ${db_file}