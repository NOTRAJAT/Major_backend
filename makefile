build:
	@go build -o bin/ecom cmd/main.go
	

run:build
	@./bin/ecom

test: 
	@go test -v ./...

migrations:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

migrate-force:
	@migrate -database -ext ext sql -dir ./cmd/migrate/migrations force  $(filter-out $@,$(MAKECMDGOALS))
