
\c postgres

create user reader with password 'reader';

create user replicator with replication encrypted password 'replicator';
select pg_create_physical_replication_slot('replication_slot');

