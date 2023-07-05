script_dir="$(cd "$(dirname "$0")" && pwd)"
cd ${script_dir}

docker build -t react:v0.1 .