all: protobuf

protobuf:
	protoc -I ./ rpc.proto --go_out=plugins=grpc:.

