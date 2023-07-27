run:
	@go run ./cmd/cmd.go 
install:
	@go build -o ./build/migrator ./cmd/cmd.go 
run-test:
	cd ./test && go test
