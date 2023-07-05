code_dir="$(cd "$(dirname "$0")" && cd .. && pwd)"

docker run -it --net=host --ipc=host --privileged=true -v ${code_dir}:/mnt --name frontend react:v0.1 /bin/bash