B=\033[0;1m
G=\033[0;92m
R=\033[0m

DIR = ${CURDIR}

.PHONY: help up down
# Show this help prompt
help:
	@echo '  Usage:'
	@echo ''
	@echo '    make <target>'
	@echo ''
	@echo '  Targets:'
	@echo ''
	@awk '/^#/{ comment = substr($$0,3) } comment && /^[a-zA-Z][a-zA-Z0-9_-]+ ?:/{ print "   ", $$1, comment }' $(MAKEFILE_LIST) | column -t -s ':' | grep -v 'IGNORE' | sort | uniq

# Stops container
down:
	@echo "\n${B}${G}Stop container${R}\n"
	@docker stop neo-raw-generator || true
	@docker rm neo-raw-generator || true

# Starts container
up: down
	@echo "\n${B}${G}build container${R}\n"
	@time docker build -t neo-bench-rpcgen .
	@echo "\n${B}${G}enter inside container:${R}\n"
	@time docker run -v ${DIR}/raw:/var/test --rm --net host -it --name neo-raw-generator neo-bench-rpcgen:latest /bin/sh
