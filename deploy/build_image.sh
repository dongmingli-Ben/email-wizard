script_dir="$(cd "$(dirname "$0")" && pwd)"
cd ${script_dir}

docker build -t nginx:v0.0 .