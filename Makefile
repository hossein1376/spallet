.PHONY: all build
all: build
build:
	@ go build -o spallet ./cmd/spallet

PHONY: test unit-test integration
test: unit-test integration
unit-test:
	@ go test -v ./pkg/application/service/...
integration:
	@ go test -v ./internal/integration/...

.PHONY: stress-load
stress-load:
	@ go run ./internal/stress

.PHONY: mock install-mockery
mock:
	@ mockery
install-mockery:
	go install github.com/vektra/mockery/v3@v3.5.4

.PHONY: gen-docs
gen-docs:
	@ node ./assets/docs/openapi/gen.js ./assets/docs/openapi/swagger.json ./assets/docs/openapi/swagger.html