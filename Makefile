setup:
# copies example config files, skipping if they already exist
	cp -n config.example.yaml secrets/config.yaml
	cp -n config.example.yaml secrets/config_test.yaml

.PHONY: run
run:
	go run main.go

.PHONY: reset_db
reset_db:
	go run api/main.go db migrate reset=true

.PHONY: migrate
migrate:
	go run api/main.go db migrate

.PHONY: generate
generate:
# Not-very-elegant way of invoking wire generator, then cleaning up the modfile:
# Because wire has dependencies that the modfile treats as unused, the modfile
# needs to be updated before and after the command
	go run -mod=mod github.com/google/wire/cmd/wire ./...
	go mod tidy

.PHONY: test
test:
	go test ./...