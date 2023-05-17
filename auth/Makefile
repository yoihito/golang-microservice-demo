docker-build: 
	docker build -t microservice-auth .

migrate-create:
	migrate create -ext sql -dir migrations -seq 'migrate_name'

bin-deps:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest