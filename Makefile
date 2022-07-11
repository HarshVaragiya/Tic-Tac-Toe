grpc:
	cd structures && protoc -I .  data.proto  --go_out=. --go_grpc_out=require_unimplemented_servers=false:.

cli:
	GOOS=darwin GOARCH=arm64 go build -o bin/m1 .
	GOOS=linux go build -o bin/linux .
	go build -o bin/tictactoe .

local: cli
	cp bin/tictactoe client1/
	cp bin/tictactoe client2/
	bin/tictactoe server
