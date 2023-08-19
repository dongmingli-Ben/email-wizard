# protoc --proto_path=../protos \
#     --go_out=. \
#     --go_opt=Mdatabase.proto=/grpc \
#     ../protos/database.proto

protoc --proto_path=../protos \
    --go_out=Mdatabase.proto=/grpc:. \
    --go-grpc_out=Mdatabase.proto=/grpc:. \
    ../protos/database.proto