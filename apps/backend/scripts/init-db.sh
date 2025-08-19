#!/bin/bash
set -e

# Create additional databases for testing
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE app_backend_test;
    GRANT ALL PRIVILEGES ON DATABASE app_backend_test TO $POSTGRES_USER;
EOSQL

echo "Additional databases created successfully"