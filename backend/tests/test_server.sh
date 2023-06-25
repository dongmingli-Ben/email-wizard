#!/bin/bash

script_dir="$(cd "$(dirname "$0")" && pwd)"

curl -G -d "user_id=toymaker&secret=xxx" http://localhost:8080/events
