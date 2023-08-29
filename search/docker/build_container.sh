set -e
cd "$(dirname "$0")"

docker network create elastic
docker run -p 9200:9200 -e ES_JAVA_OPTS="-Xms128m -Xmx128m" \
    --net elastic \
    --name search elasearch:v0.0
# docker run --name kibana --net elastic -p 5601:5601 kibana:v0.0