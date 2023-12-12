#!/usr/bin/make

start:
	docker stack deploy -c docker-compose.yaml girl

clear:
	docker stack remove girl

logs:
	docker service logs --follow girl_bot-girl

