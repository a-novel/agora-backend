COVER_FILE=$(CURDIR)/coverage.out
BIN_DIR=$(CURDIR)/bin

PKG="github.com/a-novel/agora-backend"

PKG_LIST=$(shell go list $(PKG)/... | grep -v /vendor/)

# Runs the test suite.
test:
	POSTGRES_URL=$(POSTGRES_URL_TEST) ENV="test" \
		gotestsum --packages="./..." --junitfile report.xml --format pkgname -- -count=1 -p 1 -v -coverpkg=./...

# Runs the test suite in race mode.
race:
	POSTGRES_URL=$(POSTGRES_URL_TEST) ENV="test" \
		gotestsum --packages="./..." --format pkgname -- -race -count=1 -p 1 -v -coverpkg=./...

# Run the test suite in memory-sanitizing mode. This mode only works on some Linux instances, so it is only suitable
# for CI environment.
msan:
	POSTGRES_URL=$(POSTGRES_URL_TEST) ENV="test" \
		env CC=clang env CXX=clang++ gotestsum --packages="./..." --format testname -- -msan -short $(PKG_LIST) -p 1

# Setup local environment. You might need to run `docker compose up -d` alone sometimes.
setup:
	go run ./cmd/setup/main.go
	docker compose up -d

rollback:
	go run ./cmd/rollback/main.go

rollback-test:
	go run ./cmd/rollback/main.go -d $(POSTGRES_URL_TEST)

# Starts the development server.
run:
	docker compose up -d
	go run ./cmd/main/main.go

# Plugs into the development database.
db:
	psql -h localhost -p 5432 -U postgres agora

# Plugs into the test database.
db-test:
	psql -h localhost -p 5432 -U test agora_test

# Manually rotates development JWKs. Server must be running on localhost.
rotate-keys:
	curl -X POST http://localhost:2048/api/secrets

# Print a new secret key data, using the same go standard library used in the application, to use as a mock under a
# test environment.
generate-test-key:
	go run ./cmd/utils/keys/main.go

.PHONY: all test race msan setup run db db-test rotate-keys
