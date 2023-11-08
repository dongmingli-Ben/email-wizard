git config --global --add safe.directory /mnt
script_dir="$(cd "$(dirname "$0")" && pwd)"
cd ${script_dir}

# prepare certificates
cp -r ../deploy/cert ./ 

go build -tags musl -o main .
./main

# for development (to stop the container from exiting)
sleep infinity