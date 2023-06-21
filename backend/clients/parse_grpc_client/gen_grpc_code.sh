protoc --go_out=. -I../../../protos --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ../../../protos/parse.proto