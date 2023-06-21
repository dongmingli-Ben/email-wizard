script_dir="$(cd "$(dirname "$0")" && pwd)"
cd ${script_dir}

ghz --insecure --proto ../../protos/parse.proto \
    -n 200 \
    --timeout=100s \
    --call parse.Parser.ParseEmail \
    -D example_input.json \
    127.0.0.1:50052