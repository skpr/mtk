# @todo, Needs to use skpr/mtk-mysql:latest
ARG from=previousnext/mysql:latest
FROM ${from}

ARG db_name=local
ARG db_user=local
ARG db_pass=local
ARG db_file=/workspace/db.sql

RUN database-import ${db_name} ${db_user} ${db_pass} ${db_file}
