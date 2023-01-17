RED := \033[0;91m
MAGENTA := \033[0;95m
BLUE := \033[0;94m
YELLOW := \033[93m
GREEN := \033[0;92m
NC := \033[0m

RUN = echo "${BLUE}>>> Running ${@}...${NC}"
DONE = echo "${GREEN}>>> Done ${@} without errors${NC}\n"

.PHONY: lint
lint:
	@${RUN}
	@golangci-lint run ./...
	@${DONE}

.PHONY: test
test:
	@${RUN}
	@go test -v -race -count=1 -timeout=3m
	@${DONE}
