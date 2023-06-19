cwd=$(pwd)

export PYTHONPATH=$PYTHONPATH:${cwd}/service
python email_server.py
