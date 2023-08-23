set -e

docker exec -ti search /usr/share/elasticsearch/bin/elasticsearch-reset-password -u elastic
docker exec -ti search /usr/share/elasticsearch/bin/elasticsearch-create-enrollment-token -s kibana

docker cp search:/usr/share/elasticsearch/config/certs/http_ca.crt ./config

curl --cacert config/http_ca.crt -u elastic https://localhost:9200