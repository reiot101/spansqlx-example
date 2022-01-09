package packet

//go:generate protoc -I=. --go_out . --go_opt paths=source_relative --go-grpc_out . --go-grpc_opt paths=source_relative,require_unimplemented_servers=false packet.proto
