#!/bin/bash

script_dir="$(cd "$(dirname "$0")" && pwd)"
cd ${script_dir}/tests

echo "testing email gRPC client ..."
go run test_email_client.go