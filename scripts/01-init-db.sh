#!/bin/bash
set -e

# Create two databases
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    CREATE DATABASE sourcedb;
    CREATE DATABASE targetdb;
EOSQL