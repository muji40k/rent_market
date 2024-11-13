#! /bin/bash

if [ 0 -eq $(ls -A /var/lib/postgresql/data | wc -l) ]; then
    echo 'Start Replica initialisation'
    echo 'Waiting for primary to connect...'

    until pg_basebackup --pgdata=/var/lib/postgresql/data -R \
        --slot=replication_slot --host=postgresql_db --port=5432
    do
        echo 'Waiting for primary to connect...'
        sleep 1s
    done

    echo 'Backup done, starting replica...'
    chmod 0700 /var/lib/postgresql/data
    chown postgres:root /var/lib/postgresql/data
    chown postgres:postgres -R /var/lib/postgresql/data/*
fi

su postgres - -c 'postgres'

