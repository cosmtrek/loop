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

ci: setup check
