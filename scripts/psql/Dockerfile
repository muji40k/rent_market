
FROM postgres:latest

RUN mkdir -p /docker-entrypoint-initdb.d/
RUN mkdir -p /scripts/

ADD ./config/pg_hba.conf /etc/postgresql/pg_hba.conf
ADD ./startup/*.sql /docker-entrypoint-initdb.d/
ADD ./*.sql /scripts/

