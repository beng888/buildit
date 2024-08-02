build:
	@tailwindcss -i views/css/styles.css -o public/styles.css
	@templ generate view
	@go build -o bin/buildit cmd/main.go 

test:
	@go test -v ./...
	
run: build
	@air

tailwind:
	@tailwindcss -i views/css/styles.css -o public/styles.css --watch

templ:
	@templ generate --watch --proxy=http://localhost:8080 --open-browser=false

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

# Create DB container
docker-up:
	@if docker compose up 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi