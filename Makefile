include .env

POSTGRESQL_URL ?= postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable

build:
	go build -o .out/tink .
migrate-up:
	migrate -database ${POSTGRESQL_URL} -path migrations up
migrate-down:
	migrate -database ${POSTGRESQL_URL} -path migrations down