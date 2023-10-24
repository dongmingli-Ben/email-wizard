# Main Backend Control

## Environment

First build the docker image:

```bash
cd backend
bash build_image.sh
```

Then run a container with the built image:

```bash
cd backend
bash build_container.sh
```

Attach into the container:

```bash
docker attach backend
```

## How to Run

### The Main Backend API Server

Run it with:

```bash
cd backend
bash run_server.sh
```

### Testing the Backend Server

To test whether the API endpoint is working properly, use

```bash
cd backend
bash tests/test_server.sh
```

If the API endpoint is up and healthy, you should see something like (it is slow, please be patient):

```json
[
  {
    "end_time": "2023-04-06T12:00:00 Asia/Shanghai",
    "event_type": "registration",
    "summary": "2023大学......",
    "venue": "https://....../vm/YVgulbu.aspx"
  }
]
```

To test the gRPC microservice clients, use

```bash
cd backend
bash tests/test_clients.sh
```
