# This task execute go tests with race detector and coverage in verbose mode
.PHONY: test
test:
	go test -v -cover -race .

# This task examines Go source code and reports suspicious constructs
.PHONY: vet
vet:
	go vet ./...
