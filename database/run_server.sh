git config --global --add safe.directory /mnt
cd "$(dirname "$0")"

go build -o main .
./main