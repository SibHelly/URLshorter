.PHONY: up down build rebuild logs clean ps db-clean

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

# Очистка хранилища базы данных (удаляет volume)
db-clean:
	docker volume rm urlshorter_postgres_data

# Показать список контейнеров
ps:
	docker-compose ps