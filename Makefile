TSURU_DEPLOY_FILES 	 			:= main

TSURU_APP_NAME  		:= check-password-prod
LOG_LEVEL				:= info
PROXY_URL				:= http://proxy.globoi.com:3128

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

######################## TSURU ########################
# Set environment variables
.PHONY: set-tsuru-env-variables-prod
set-tsuru-env-variables-prod:
	tsuru env-set LOG_LEVEL=$(LOG_LEVEL) PROXY_URL=$(PROXY_URL) -a $(TSURU_APP_NAME) --no-restart

# Deploy to tsuru
.PHONY: deploy-tsuru-prod
deploy-tsuru-prod: set-tsuru-env-variables-prod
	GOOS=linux GOARCH=amd64 go build cmd/main.go
	tsuru app-deploy public static $(TSURU_DEPLOY_FILES) Procfile tsuru.yaml -a $(TSURU_APP_NAME)
