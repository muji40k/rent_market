
FROM postgres:latest

RUN mkdir -p /docker-entrypoint-initdb.d/
RUN mkdir -p /scripts/

ADD ./*.sql /scripts/
ADD ./startup/*.sql /docker-entrypoint-initdb.d/

