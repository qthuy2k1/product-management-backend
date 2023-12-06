run: 
	docker compose up
  
test:
	go test ./... -v

postgres:
	psql -U postgres -d product_management -h localhost -p 5432

create_migration:
	migrate create -ext sql -dir data/migrations/ -seq $(filename)

# Change %cd% to $(pwd) if your terminal is not Command Prompt
migrate_up:
	docker run --rm -it --network product-management_host -v "%cd%/data/migrations:/migrations" migrate/migrate -path=/migrations/  -database "postgres://postgres:root@host.docker.internal:5432/product_management?sslmode=disable" up $(ver)

# Change %cd% to $(pwd) if your terminal is not Command Prompt
migrate_up_force:
	docker run --rm -it --network product-management_host -v "%cd%/data/migrations:/migrations" migrate/migrate -path=/migrations/  -database "postgres://postgres:root@host.docker.internal:5432/product_management?sslmode=disable" force $(ver)

# Change %cd% to $(pwd) if your terminal is not Command Prompt
migrate_down:
	docker run --rm -it --network product-management_host -v "%cd%/data/migrations:/migrations" migrate/migrate -path=/migrations/  -database "postgres://postgres:root@host.docker.internal:5432/product_management?sslmode=disable" down $(ver)

gen-mocks:
	mockery --dir internal/controllers --all --recursive --inpackage
	mockery --dir internal/repositories --all --recursive --inpackage

gen-graph:
	go run github.com/99designs/gqlgen generate
