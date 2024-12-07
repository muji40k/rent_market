
FROM postgres:latest

RUN mkdir -p /docker-entrypoint-initdb.d/
RUN mkdir -p /scripts/
RUN mkdir -p /testdata/

ADD ./benchmark/*.sql /docker-entrypoint-initdb.d/
ADD ./*.sql /scripts/
ADD ./inserts/categories.sql /testdata/
ADD ./inserts/bench/* /testdata/

