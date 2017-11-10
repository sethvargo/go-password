test:
	@go test ./...
.PHONY: test

test-race:
	@go test -race ./...
.PHONY: test-race
