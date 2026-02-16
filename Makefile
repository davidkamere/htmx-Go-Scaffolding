APP_BIN=./tmp/server

.PHONY: deps css test vet fmt fmt-check run build dev ci

deps:
	npm install

css:
	npm run build:css
	npm run build:vendor

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w $$(git ls-files '*.go')

fmt-check:
	@unformatted="$$(gofmt -l $$(git ls-files '*.go'))"; \
	if [ -n "$$unformatted" ]; then \
		echo "Files not formatted:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

run:
	go run ./cmd/server

build:
	go build -o $(APP_BIN) ./cmd/server

dev:
	air

ci:
	npm ci
	npm run build:css
	npm run build:vendor
	go test ./...
	go vet ./...
	$(MAKE) fmt-check
