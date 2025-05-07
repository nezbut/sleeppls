build:
	go build -o bin/cli/sleeppls-cli cmd/sleeppls-cli/main.go

install:
	go install ./cmd/sleeppls-cli

run-one:
	go run ./cmd/sleeppls-cli

clean:
	rm -rf ./bin/

deps-download:
	@echo "Downloading Go modules..."
	go mod download

deps-tidy:
	@echo "Tidying Go modules..."
	go mod tidy

deps-update:
	@echo "Update deps"
	go get -u ./...
	$(MAKE) deps-tidy
