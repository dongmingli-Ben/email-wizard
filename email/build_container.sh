code_dir="$(cd "$(dirname "$0")" && cd .. && pwd)"

docker run -it --net=host --ipc=host --privileged=true -v ${code_dir}:/mnt --name email pyemail:v0.2 /bin/bash