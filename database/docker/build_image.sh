set -e

script_dir="$(cd "$(dirname "$0")" && pwd)"
cd ${script_dir}

# docker build -f Dockerfile.postgresql -t postgresql:v0.0 .
docker build -f Dockerfile.golang -t gopostgre:v0.1 .