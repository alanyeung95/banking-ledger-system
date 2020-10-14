.DEFAULT_GOAL := run

NETWORK := demo.network
NETWORK_EXIST = $(shell docker network ls | grep -E ''"$(NETWORK)"'([^-]|$$)' | wc -l | tr -d ' \t\n\r\f')

.PHONY: network
network:
ifeq ($(NETWORK_EXIST), 0)
	@echo "Building docker network $(NETWORK)"
	docker network create $(NETWORK)
endif

.PHONY: run
run: network
	docker-compose up

.PHONY: test
test: network
	docker-compose -f docker-compose.yaml -f docker-compose.test.yaml up --abort-on-container-exit

.PHONY: clean
clean: 
	docker-compose down --remove-orphans