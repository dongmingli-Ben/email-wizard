cd "$(dirname "$0")"

docker build -f Dockerfile.elastic -t elasearch:v0.0 .
# docker build -f Dockerfile.kibana -t kibana:v0.0 .