all: build
build:
	@ go build -o spallet ./cmd/spallet

.PHONY: gen-docs
gen-docs:
	@ node ./assets/docs/openapi/gen.js ./assets/docs/openapi/swagger.json ./assets/docs/openapi/swagger.html