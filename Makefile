.PHONY: default
default: clear-port-forwards kind-check
	skaffold run --port-forward -f skaffold.yaml &

.PHONY: debug
debug: clear-port-forwards kind-check
	skaffold run --profile=debug --port-forward -f skaffold.yaml &

.PHONY: sync
sync: clear-port-forwards kind-check
	skaffold dev -f skaffold.yaml

.PHONY: sync-debug
sync-debug: clear-port-forwards kind-check
	skaffold dev --profile=debug -f skaffold.yaml

.PHONY: kind-check
kind-check:
	@if [ "$(shell kubectl config current-context)" != "kind-kind" ]; then \
		kind create cluster --name kind; \
		kubectl config use-context kind-kind; \
	fi

.PHONY: test
test: kind-check
	helm test startupers
	skaffold verify -f skaffold.yaml
	skaffold verify -f skaffold.yaml
	go test ./...

.PHONY: benchmark
benchmark: kind-check
	go test -bench . ./...
	k6 run internal/k6s/get-users.js
	# k6 run --env CONNECTION_STRING="http://localhost:9000" internal/k6s/get-users.js

.PHONY: snyk
snyk:
	@command -v snyk >/dev/null 2>&1 || { echo >&2 "Snyk CLI is not installed. Please install it from https://snyk.io/"; exit 1; }
	snyk test

.PHONY: clean
clean: clear-port-forwards
	kind delete cluster
	docker system prune

.PHONY: render
render:
	skaffold render

.PHONY: docs
docs:
	swag init -g cmd/main.go

.PHONY: clear-port-forwards
clear-port-forwards:
	-killall skaffold

.PHONY: seed
seed: kind-check
	PGPASSWORD=password psql -h localhost -U username -d startupers -c "INSERT INTO users (name) VALUES ('Alice') ON CONFLICT DO NOTHING;"
	PGPASSWORD=password psql -h localhost -U username -d startupers -c "INSERT INTO users (name) VALUES ('Bob') ON CONFLICT DO NOTHING;"
	PGPASSWORD=password psql -h localhost -U username -d startupers -c "INSERT INTO users (name) VALUES ('Charlie') ON CONFLICT DO NOTHING;"

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  default    - Run skaffold with default settings"
	@echo "  debug      - Run skaffold in debug mode"
	@echo "  sync       - Run skaffold in dev mode"
	@echo "  sync-debug - Run skaffold in dev mode with debug profile"
	@echo "  kind-check - Check if the current Kubernetes context is 'kind', create a 'kind' cluster if not"
	@echo "  test       - Run Helm tests, skaffold verify, and Go tests"
	@echo "  benchmark  - Run Go benchmarks and k6 performance tests"
	@echo "  snyk       - Run Snyk tests (checks if the Snyk CLI is installed)"
	@echo "  clean      - Delete the 'kind' cluster"
	@echo "  seed       - Add test data to database"
	@echo "  render     - Run skaffold render"
	@echo "  docs       - Generate Swagger documentation"
