cwd=$(pwd)

export PYTHONPATH=$PYTHONPATH:${cwd}/service
python parse_server.py
