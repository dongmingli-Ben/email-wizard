code_dir="$(cd "$(dirname "$0")" && cd .. && pwd)"

docker run -it --net=host --ipc=host --privileged=true -v ${code_dir}:/mnt --name backend goback:v0.3 /bin/bash