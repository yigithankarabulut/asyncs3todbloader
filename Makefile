DC := docker-compose -f ./docker-compose.yml
PORT := $(word 2,$(MAKECMDGOALS))
METHOD := $(word 3,$(MAKECMDGOALS))
PRODUCT_ID := $(word 4,$(MAKECMDGOALS))
HOST := "http://localhost:"
URI := "/api/v1/product/"

all:
	@$(DC) up -d

build:
	@$(DC) up --build

down:
	@$(DC) down

re: down all

clean:
	@$(DC) down -v

db:
	@$(DC) up -d --build mongodb

job:
	@$(DC) up --build job

microservice:
	@$(DC) up --build microservice

request:
	@curl -X "$(METHOD)" "$(HOST)$(PORT)$(URI)$(PRODUCT_ID)"
	@echo ""
%:
	@:

.PHONY: all build down re clean db job microservice request