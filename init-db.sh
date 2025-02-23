#!/bin/bash
set -e

# Create the liivedb database if it doesn't exist
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-'EOSQL'
    CREATE DATABASE liivedb;
    GRANT ALL PRIVILEGES ON DATABASE liivedb TO liive;
EOSQL 