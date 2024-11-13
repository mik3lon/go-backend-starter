ifneq (,$(wildcard .env))
  include .env
  export
endif

# Variables
DOCKER_COMPOSE_FILE = docker-compose.yml
YELLOW  = \033[33m
CYAN    = \033[36m
GREEN   = \033[32m
RESET   = \033[0m

# Default target: show help message
.PHONY: help up down logs clean

help:
	@echo ""
	@echo "Usage: ${CYAN}make [target]${RESET}"
	@echo ""
	@echo "${GREEN}Available targets:${RESET}"
	@echo ""
	@echo "  ${YELLOW}up${RESET}      - Start the application containers in detached mode"
	@echo "  ${YELLOW}down${RESET}    - Stop the application containers and remove the network"
	@echo "  ${YELLOW}logs${RESET}    - View logs of the running containers in real-time"
	@echo "  ${YELLOW}clean${RESET}   - Stop the containers and remove volumes and orphan containers"
	@echo ""
	@echo "${CYAN}Example: make up${RESET}"
	@echo ""

# Start the application containers
up:
	@echo "${CYAN}Starting application containers...${RESET}"
	docker compose -f $(DOCKER_COMPOSE_FILE) up -d

# Stop the application containers
down:
	@echo "${CYAN}Stopping application containers...${RESET}"
	docker compose -f $(DOCKER_COMPOSE_FILE) down

# View logs of the running containers
logs:
	@echo "${CYAN}Showing logs from running containers...${RESET}"
	docker compose -f $(DOCKER_COMPOSE_FILE) logs -f

# Remove all containers, networks, and volumes created by docker compose
clean:
	@echo "${CYAN}Cleaning up containers, volumes, and networks...${RESET}"
	docker compose -f $(DOCKER_COMPOSE_FILE) down --volumes --remove-orphans

# execute tests
tests:
	@echo "${CYAN} Running tests... ${RESET}"
	@go test ./... \
	&& echo "${GREEN}Success: All tests passed!${RESET}" \
	|| echo "${RED}Error: Some tests failed.${RESET}"
