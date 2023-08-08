code_dir="$(cd "$(dirname "$0")" && cd .. && pwd)"

docker run -it --rm -d --net=host --ipc=host --privileged=true -v ${code_dir}/frontend/client/dist:/usr/share/nginx/html --name nginx nginx:v0.0