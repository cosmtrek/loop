LDFLAGS += -X "main.BuildTimestamp=$(shell date -u "+%Y-%m-%d %I:%M:%S")"
LDFLAGS += -X "main.Version=$(shell git rev-parse HEAD)"

setup:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/golang/lint/golint
	go get -u github.com/Masterminds/glide
	glide install
	@echo "Install pre-commit hook"
	ln -s $(shell pwd)/hooks/pre-commit $(shell pwd)/.git/hooks/pre-commit

.PHONY: check
check:
	@./hack/check.sh ${scope}

.PHONY: release
release: check
	@mkdir -p bin
	GOOS=linux go build -ldflags '$(LDFLAGS)' -o bin/loop-linux
	GOOS=darwin go build -ldflags '$(LDFLAGS)' -o bin/loop-darwin

ci: setup check
