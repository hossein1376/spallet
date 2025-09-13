.PHONY: all build run
all: build
build:
	@ go build -o spallet ./cmd/spallet
run:
	@ go run ./cmd/spallet

PHONY: test unit-test integration
test: unit-test integration
unit-test:
	@ go test ./pkg/application/service/... -v
integration:
	@ go test -C ./internal/test/integration/ -v

.PHONY: stress-load
stress-load:
	@ go run ./internal/stress

.PHONY: mock install-mockery
mock:
	@ mockery
install-mockery:
	@ go install github.com/vektra/mockery/v3@v3.5.4

.PHONY: gen-docs
gen-docs:
	@ node ./assets/docs/openapi/gen.js ./assets/docs/openapi/swagger.json \
		./assets/docs/openapi/swagger.html