script_dir="$(cd "$(dirname "$0")" && pwd)"
cd ${script_dir}

export PYTHONPATH=$PYTHONPATH:${script_dir}/service
python email_server.py & python email_server_async.py

sleep infinity