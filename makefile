migrate:
	go run cmd/cmd.go migrate
rollback:
	go run cmd/cmd.go rollback
install:
	go build -o ./build/migrator ./cmd/cmd.go 

