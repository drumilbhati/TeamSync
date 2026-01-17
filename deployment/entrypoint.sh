#!/bin/bash
set -e

# Setup Postgres
export PGDATA=/var/lib/postgresql/data
export PGBIN=/usr/lib/postgresql/14/bin

if [ ! -d "$PGDATA/base" ]; then
    echo "Initializing database..."
    mkdir -p "$PGDATA"
    chown -R postgres:postgres "$PGDATA"
    # Ensure the directory is empty for initdb
    rm -rf "$PGDATA"/*
    su postgres -c "$PGBIN/initdb -D $PGDATA"
    
    # Start postgres temporarily to create user and schema
    su postgres -c "$PGBIN/pg_ctl -D $PGDATA -l /tmp/pg_log start"
    
    echo "Waiting for postgres to start..."
    # Loop until pg_isready returns 0
    MAX_RETRIES=30
    COUNT=0
    while ! su postgres -c "$PGBIN/pg_isready" && [ $COUNT -lt $MAX_RETRIES ]; do
      echo "Waiting... ($COUNT)"
      sleep 1
      ((COUNT++))
    done

    if [ $COUNT -eq $MAX_RETRIES ]; then
      echo "Postgres failed to start"
      cat /tmp/pg_log
      exit 1
    fi

    echo "Creating user and database..."
    su postgres -c "$PGBIN/psql --command \"CREATE USER \\\"user\\\" WITH SUPERUSER PASSWORD 'password';\""
    su postgres -c "$PGBIN/createdb -O user teamsync"
    
    echo "Applying schema..."
    su postgres -c "$PGBIN/psql -d teamsync -f /app/database/schema.sql"
    
    su postgres -c "$PGBIN/pg_ctl -D $PGDATA stop"
fi

# Fix permissions for Postgres
chown -R postgres:postgres /var/lib/postgresql/data
chmod 700 /var/lib/postgresql/data

# Replace PORT in nginx config
# On Ubuntu, we should update both available and enabled if they are separate
sed -i "s/RENDER_PORT/$PORT/g" /etc/nginx/sites-available/default

# Start Supervisor
echo "Starting services via Supervisor..."
exec /usr/bin/supervisord -c /etc/supervisor/conf.d/supervisord.conf
