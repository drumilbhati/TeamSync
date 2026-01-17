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
    su - postgres -c "$PGBIN/initdb -D $PGDATA"
    
    # Start postgres temporarily to create user and schema
    su - postgres -c "$PGBIN/pg_ctl -D $PGDATA -l /tmp/pg_log start"
    
    echo "Waiting for postgres to start..."
    until su - postgres -c "$PGBIN/pg_isready"; do
      sleep 1
    done

    echo "Creating user and database..."
    su - postgres -c "psql --command \"CREATE USER \\\"user\\\" WITH SUPERUSER PASSWORD 'password';\""
    su - postgres -c "createdb -O user teamsync"
    
    echo "Applying schema..."
    su - postgres -c "psql -d teamsync -f /app/database/schema.sql"
    
    su - postgres -c "$PGBIN/pg_ctl -D $PGDATA stop"
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
