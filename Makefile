test:
	@go test \
		-count=1 \
		-race \
		-shuffle=on \
		-timeout=5m \
		-vet=all \
		./...
.PHONY: test
