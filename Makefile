# Determine if you have docker-compose or docker compose installed locally
# If this does not work on your system, just set the name of the executable you have installed
DCO_BIN := $(shell { command -v docker-compose || command -v docker compose; } 2>/dev/null)
# Handle Spaces in Docker Compose path
DCO_BIN := $(subst $(null) $(null),\ ,$(DCO_BIN))
# Connect to the primary database
.PHONY: db
db:
	docker exec -it db_couploan psql postgresql://postgres:root@localhost:5437/couploan

# Connect to the test database (you must run tests first before running this)
.PHONY: db-test
db-test:
	docker exec -it db_couploan psql postgresql://postgres:root@localhost:5437/couploan_test

# Connect to the primary cache
.PHONY: cache
cache:
	docker exec -it cache_couploan redis-cli

 # Connect to the test cache
.PHONY: cache-test
cache-test:
	docker exec -it cache_couploan redis-cli -n 1

# Install Ent code-generation module
.PHONY: ent-install
ent-install:
	go get -d entgo.io/ent/cmd/ent

# Generate Ent code
.PHONY: ent-gen
ent-gen:
	go generate ./ent

# Create a new Ent entity
.PHONY: ent-new
ent-new:
	go run entgo.io/ent/cmd/ent new $(name)

# Start the Docker containers
.PHONY: up
up:
	$(DCO_BIN) up -d
	sleep 3

# Rebuild Docker containers to wipe all data
.PHONY: reset
reset:
	$(DCO_BIN) down
	make up

# Run the application
.PHONY: watch
watch:
	clear
	air

# Build the application
.PHONY: build
build:
	clear
	go build -o ./tmp/server.exe ./cmd/server/

# Run the application
.PHONY: run
run:
	clear
	go run ./cmd/server/

# Run all tests
.PHONY: test
test:
	go test -count=1 -p 1 ./...

# Check for direct dependency updates
.PHONY: check-updates
check-updates:
	go list -u -m -f '{{if not .Indirect}}{{.}}{{end}}' all | grep "\["

# Run the asynqmon
.PHONY: asynqmon
asynqmon:
	docker run --rm \
        --name couploan_asynqmon \
        -p 9990:9990 \
        hibiken/asynqmon --redis-url=redis://host.docker.internal:6384/7 --port=9990

.PHONY: atlas-down
atlas-down:
	atlas migrate down -c file://atlas.hcl --dir file://ent/migrate/migrations --url "postgres://postgres:root@localhost:5437/couploan?search_path=public&sslmode=disable"

.PHONY: atlas-down-test
atlas-down-test:
	atlas migrate down -c file://atlas.hcl --dir file://ent/migrate/migrations --url "postgres://postgres:root@localhost:5437/couploan_test?search_path=public&sslmode=disable"

.PHONY: atlas-apply
atlas-apply: atlas-validate
	atlas migrate apply -c file://atlas.hcl --dir file://ent/migrate/migrations --url "postgres://postgres:root@localhost:5437/couploan?search_path=public&sslmode=disable"

.PHONY: atlas-apply-test
atlas-apply-test: atlas-validate
	atlas migrate apply -c file://atlas.hcl --dir file://ent/migrate/migrations --url "postgres://postgres:root@localhost:5437/couploan_test?search_path=public&sslmode=disable"

.PHONY: atlas-diff
atlas-diff: atlas-validate ent-gen
	 go run ent/migrate/main.go $(name)

.PHONY: atlas-diff-test
atlas-diff-test: atlas-validate ent-gen atlas-clean-test
	go run ent/migrate/main.go $(name) test

.PHONY: atlas-status
atlas-status:
	atlas migrate status -c file://atlas.hcl --url "postgres://postgres:root@localhost:5437/couploan?search_path=public&sslmode=disable" --dir file://ent/migrate/migrations

.PHONY: atlas-status-test
atlas-status-test:
	atlas migrate status -c file://atlas.hcl --url "postgres://postgres:root@localhost:5437/couploan_test?search_path=public&sslmode=disable" --dir file://ent/migrate/migrations

.PHONY: atlas-validate
atlas-validate: atlas-hash
	atlas migrate validate -c file://atlas.hcl --dir file://ent/migrate/migrations

.PHONY: atlas-hash
atlas-hash:
	atlas migrate hash -c file://atlas.hcl --dir file://ent/migrate/migrations

.PHONY: atlas-new
atlas-new:
	atlas migrate new -c file://atlas.hcl $(name) --dir file://ent/migrate/migrations

.PHONY: atlas-clean
atlas-clean:
	atlas schema clean $(name) --url "postgres://postgres:root@localhost:5437/couploan?search_path=public&sslmode=disable"

.PHONY: atlas-clean-test
atlas-clean-test:
	atlas schema clean $(name) --url "postgres://postgres:root@localhost:5437/couploan_test?search_path=public&sslmode=disable"

.PHONY: release
release: atlas-validate
	goreleaser release --clean --snapshot