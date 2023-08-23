set -e
cd "$(dirname "$0")"

password="$(cat config/config.json | jq -r '.password')"
curl -X DELETE --cacert config/http_ca.crt -u "elastic:${password}" https://localhost:9200/email-wizard-events