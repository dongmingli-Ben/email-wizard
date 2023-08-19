cd "$(dirname "$0")"

protoc --proto_path=../protos \
    --go_out=Mdatabase.proto=/grpc:. \
    --go-grpc_out=Mdatabase.proto=/grpc:. \
    ../protos/database.proto