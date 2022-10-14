generate:
	protoc --go_out=protos --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative protos/moviesapp.proto

run-server:
	go run server/main.go

run-client:
	go run client/main.go