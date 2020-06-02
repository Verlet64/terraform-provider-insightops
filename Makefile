BIN=cmd/terraform-provider-insightops

.PHONY: test
test: 
	go test ./...

.PHONY: build
build:
	cd $(BIN) && go build