export COMPOSE_DOCKER_CLI_BUILD=1
export DOCKER_BUILDKIT=1

all_d: down build up_d build_db

all: down build up

up:
	docker-compose up

up_d:
	docker-compose up -d

build:
	docker-compose build

down:
	docker-compose down --remove-orphans

build_db:
	psql postgres://postgres:postgres@localhost:5432/postgres < ./sql/create_tables.sql

