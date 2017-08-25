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

.PHONY: run
run: docker-compose-build-api docker-compose-up-api

.PHONY: build-hashcar-runner
build-hashcar-runner:
	GOOS=linux GOARCH=amd64 go build cmd/hashcat_runner.go
