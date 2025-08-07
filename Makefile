DBURL=postgres://postgres:postgres@localhost:5432/affarm
mgtbl=migrations

gen-proto:
	protoc -I proto proto/${pkg}/*.proto --go_out=proto/gen/ --go_opt=paths=source_relative --go-grpc_out=proto/gen/ --go-grpc_opt=paths=source_relative

mock-services:
	cd ./internal/services && mockery --all

gen-docs:
	swag init -g=internal/controller/rest/v1/router.go -o=docs --parseInternal

migration-up:
	migrate -path migrations -database '${DBURL}?x-migrations-table=${mgtbl}&sslmode=disable' up ${version}

migration-down:
	migrate -path migrations -database '${DBURL}?x-migrations-table=${mgtbl}&sslmode=disable' down ${version}

migration-force:
	migrate -path migrations -database "${DBURL}?x-migrations-table=${mgtbl}&sslmode=disable" force ${version}

run:
	go run cmd/app/main.go --config=config/local.env