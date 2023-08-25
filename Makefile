.DEFAULT_GOAL := help
# Name of the aws profile you want to use
AWS_PROFILE = "medallion"

#file use to deploy in a easiest way

GOOS=linux
GOARCH=amd64
CGO_ENABLED=0


.PHONY: help
help: # Show all the commands
		@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

.PHONY: run
run: # Run the application in local using docket
	AWS_PROFILE=${AWS_PROFILE} sam local start-api

.PHONY: build
build: # Build the packages of go and create the template output for deployment
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 sam build

.PHONY: build-run
build-run: # Option to make a build and run local the microservices
	make build && make run

.PHONY: first-deploy
first-deploy: # Option to make a first deploy
	make build && sam deploy --guided --profile ${AWS_PROFILE}

.PHONY:deploy
deploy: # Posterior deployments
	make build && sam deploy

.PHONY: delete
delete: # delete the template and all resourses in aws
	sam delete --profile ${AWS_PROFILE}