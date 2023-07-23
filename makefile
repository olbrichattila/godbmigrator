migrate:
	go run cmd/cmd.go migrate
rollback:
	go run cmd/cmd.go rollback
install:
	go build -o ./build/migrator ./cmd/cmd.go 
switch-sqlite:
	cp .env.sqlite.example .env
switch-mysql:
	cp .env.mysql.example .env
switch-pgsql:
	cp .env.pgsql.example .env
