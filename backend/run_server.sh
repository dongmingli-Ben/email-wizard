git config --global --add safe.directory /mnt
script_dir="$(cd "$(dirname "$0")" && pwd)"
cd ${script_dir}

# prepare certificates
cp -r ../deploy/cert ./ 

go build -o main .
./main