git config --global --add safe.directory /mnt
cd "$(dirname "$0")"

export PGPASSWORD=123456

go build -o main ./grpc_server
go build -tags musl -o async ./async_server
./main & ./async

sleep infinity