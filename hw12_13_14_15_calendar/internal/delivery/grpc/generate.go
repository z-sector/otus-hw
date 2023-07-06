package grpc

//go:generate protoc --proto_path=./../../../ --go_out=. --go-grpc_out=. ./../../../api/proto/event/event.proto
//go:generate protoc --proto_path=./../../../ --go_out=. --go-grpc_out=. ./../../../api/proto/internal/internal.proto
