apt update
apt install -y protobuf-compiler
cd ./proto
protoc --go_out=. --go-grpc_out=require_unimplemented_servers=false:. ./user.proto