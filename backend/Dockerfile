FROM golang:1.20.7-alpine3.18
RUN apk update && apk add git bash protobuf curl jq
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
RUN apk add --no-cache gcc musl-dev
WORKDIR /mnt