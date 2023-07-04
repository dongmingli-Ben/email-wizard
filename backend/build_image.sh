script_dir="$(cd "$(dirname "$0")" && pwd)"
cd ${script_dir}

docker build -t goback:v0.3 .