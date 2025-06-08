.PHONY: up down build rebuild logs clean ps

up:
	docker-compose up -d

down:
	docker-compose down

build:
	docker-compose build

rebuild: down build up

logs:
	docker-compose logs -f

db-logs:
	docker-compose logs -f db

db-bash:
	docker-compose exec db bash

db-psql:
	docker-compose exec db psql -U user -d postgres

url-logs:
	docker-compose logs -f url-shortener

url-bash:
	docker-compose exec url-shortener sh

