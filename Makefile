
.DEFAULT_GOAL := help

.PHONY: test
test: ## run tests (using go1.16beta1 for now)
	POSTGRES_DSN="postgres://marshalexamples:marshalexamples@127.0.0.1:6543/marshalexamples?sslmode=disable" go1.16beta1 test -v -race .

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
