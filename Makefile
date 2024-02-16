app-name = "golang-example"

.PHONY: clean install unittest build docker run stop vendor migrate

install:
	go get -v github.com/rubenv/sql-migrate/...
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	cd configs && sqlc generate

run:
	sql-migrate up -env=default-factories -config=configs/dbconfig.yml
	sql-migrate up -env=default-migrations -config=configs/dbconfig.yml
	sql-migrate up -env=default-seeds -config=configs/dbconfig.yml
	cd configs && sqlc generate
	swag init --output=generated/docs --parseDependency --parseInternal
	go run main.go

sqlc:
	cd configs && sqlc generate