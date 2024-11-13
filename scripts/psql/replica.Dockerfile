
FROM postgres:latest

ADD ./start_replica.sh /start_replica.sh

ENTRYPOINT ["/start_replica.sh"]

