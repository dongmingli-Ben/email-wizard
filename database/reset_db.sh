set -e
export PGPASSWORD=123456

psql -U postgres -d email-wizard-data -h postgres -t \
    -c "SELECT 'DROP TABLE IF EXISTS ' || table_name || ' CASCADE;' FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE';" \
    > clear_db.sql

psql -U postgres -d email-wizard-data -h postgres -f clear_db.sql
