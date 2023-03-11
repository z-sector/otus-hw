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

.PHONY: bench
bench:
	@${RUN}
	@go test -bench=BenchmarkGetDomainStat -benchmem -benchtime=10s -count=5 | tee new
	@${DONE}

.PHONY: benchstat
benchstat: bench
	@${RUN}
	@benchstat old new
	@${DONE}

.PHONY: test_bench
test_bench:
	@${RUN}
	@go test -v -count=1 -timeout=30s -tags=bench .
	@${DONE}

.PHONY: test_common
test_common:
	@${RUN}
	@go test -v -count=1 -timeout=30s .
	@${DONE}

.PHONY: test
test: test_common test_bench
