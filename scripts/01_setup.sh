#!/bin/bash
set -e

# PARAMS="--host=$DB_HOST --user=$USER --password=$PASSWORD"
# cat schema.sql | postgres $PARAMS

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER docker;
    GRANT ALL PRIVILEGES ON DATABASE openencoder TO docker;
EOSQL