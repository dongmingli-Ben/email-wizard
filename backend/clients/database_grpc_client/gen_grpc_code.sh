cd "$(dirname "$0")" && cd ../..

protoc --proto_path=../protos \
    --go_out=Mdatabase.proto=/clients/database_grpc_client:. \
    --go-grpc_out=Mdatabase.proto=/clients/database_grpc_client:. \
    ../protos/database.proto