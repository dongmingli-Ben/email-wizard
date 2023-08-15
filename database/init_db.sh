set -e

cd "$(dirname "$0")"

psql -h postgres -d email-wizard-data -U postgres -f create_database.sql