git config --global --add safe.directory /mnt
script_dir="$(cd "$(dirname "$0")" && pwd)"
cd ${script_dir}

go build -o main .
./main