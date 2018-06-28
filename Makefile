build-proto: build-proto-account

build-proto-account:
	protoc ./proto/accounts.proto -I. --go_out=plugins=grpc:$GOPATH/src