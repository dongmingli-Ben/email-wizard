docker run -it --net=host --ipc=host --privileged=true -v /home/toymaker/projects/email-wizard:/mnt --name backend goback:v0.2 /bin/bash