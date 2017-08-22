TSURU_DEPLOY_FILES 	 			:= main

# DEV
TSURU_APP_NAME_DEV  		:= check-password-dev
LOG_LEVEL_DEV				:= info

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

######################## DEV ########################
# Set environment variables
.PHONY: set-tsuru-env-variables-dev
set-tsuru-env-variables-dev:
	tsuru env-set LOG_LEVEL=$(LOG_LEVEL_DEV) -a $(TSURU_APP_NAME_DEV) --no-restart

# Deploy to tsuru
.PHONY: deploy-tsuru-dev
deploy-tsuru-dev:
	GOOS=linux GOARCH=amd64 go build cmd/main.go
	tsuru app-deploy public static $(TSURU_DEPLOY_FILES) Procfile tsuru.yaml -a $(TSURU_APP_NAME_DEV)
