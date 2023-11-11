script_dir="$(cd "$(dirname "$0")" && pwd)"
cd ${script_dir}

export PYTHONPATH=$PYTHONPATH:${script_dir}/service
source setup_openai.sh
python parse_server.py & python parse_server_async.py

sleep infinity
