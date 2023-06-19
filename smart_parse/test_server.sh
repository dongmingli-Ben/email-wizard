cwd=$(pwd)

export PYTHONPATH=$PYTHONPATH:${cwd}/service
python tests/test_server.py
