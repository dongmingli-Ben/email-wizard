set -e

psql -h postgres -d email-wizard-data -U postgres -f create_database.sql