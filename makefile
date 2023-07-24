run:
	@go run ./cmd/cmd.go 
install:
	@go build -o ./build/migrator ./cmd/cmd.go 
