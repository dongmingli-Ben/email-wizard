git config --global --add safe.directory /mnt
cd "$(dirname "$0")"

export PGPASSWORD=123456

go build -o main ./grpc_server
./main

sleep infinity