script_dir="$(cd "$(dirname "$0")" && pwd)"
cd ${script_dir}

export PYTHONPATH=$PYTHONPATH:${script_dir}/service
python parse_server.py
