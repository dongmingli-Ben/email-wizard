set -e

code_dir="$(cd "$(dirname "$0")" && cd ../.. && pwd)"

docker network create postgresql-net
docker run -d --ipc=host --privileged=true -v ${code_dir}:/mnt --name postgres \
     --network postgresql-net -e POSTGRES_PASSWORD=123456 -e POSTGRES_DB=email-wizard-data postgresql:v0.0
docker run -it --ipc=host --privileged=true -v ${code_dir}:/mnt --name data \
     --network postgresql-net -p 50100-50200:50100-50200 gopostgre:v0.1 bash
