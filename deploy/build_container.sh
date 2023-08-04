code_dir="$(cd "$(dirname "$0")" && cd .. && pwd)"

docker run -it --net=host --ipc=host --privileged=true -v ${code_dir}:/mnt --name nginx nginx:v0.0 /bin/bash