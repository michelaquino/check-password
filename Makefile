######################## LOCAL ########################
######## Docker compose ########
.PHONY:	docker-compose-build-api
docker-compose-build-api:
	docker-compose build

.PHONY:	docker-compose-up-api
docker-compose-up-api:
	docker-compose up

.PHONY:	docker-compose-stop-api
docker-compose-stop-api:
	docker-compose stop

######## Deploy ########
.PHONY: run
run: docker-compose-build-api docker-compose-up-api
