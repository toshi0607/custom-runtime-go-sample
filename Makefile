################
# COMMON
################

build: build-function build-bootstrap
.PHONY: build

build-function:
	GOARCH=amd64 GOOS=linux go build -o artifacts/handler ./handlers
.PHONY: build-function

build-bootstrap:
	GOARCH=amd64 GOOS=linux go build -o artifacts/bootstrap ./init
.PHONY: build-bootstrap



################
# SAM
################

STACK_NAME := custom-runtime-go-sample
TEMPLATE_FILE := template.yml
SAM_FILE := sam.yml

deploy: build
	sam package \
		--template-file $(TEMPLATE_FILE) \
		--s3-bucket $(STACK_BUCKET) \
		--output-template-file $(SAM_FILE)
	sam deploy \
		--template-file $(SAM_FILE) \
		--stack-name $(STACK_NAME) \
		--capabilities CAPABILITY_IAM
.PHONY: deploy

create-bucket:
	aws s3 mb "s3://$(STACK_BUCKET)"
.PHONY: create-bucket

delete:
	aws cloudformation delete-stack --stack-name $(STACK_NAME)
	aws s3 rm "s3://$(STACK_BUCKET)" --recursive
	aws s3 rb "s3://$(STACK_BUCKET)"
.PHONY: delete



################
# AWS CLI
################

FUNCTION_NAME := go-custome-runtime-sample
ARTIFACT_ZIP := artifacts.zip

create-zip:
	cp artifacts/handler handler
	cp artifacts/bootstrap bootstrap
	chmod +x handler bootstrap
	zip $(ARTIFACT_ZIP) handler bootstrap
	rm handler
	rm bootstrap
.PHONY: create-zip

deploy-sub: build create-zip update-function
.PHONY: deploy-sub

recreate-function: delete-function create-function
.PHONY: recreate-function

create-function: build create-zip
	aws lambda create-function \
	  --function-name $(FUNCTION_NAME) \
	  --zip-file "fileb://$(ARTIFACT_ZIP)" \
	  --handler "handler" \
	  --runtime provided \
	  --role $(LAMBDA_ROLL_ARN)
.PHONY: create-function

delete-function:
	aws lambda delete-function \
	  --function-name $(FUNCTION_NAME)
.PHONY: delete-function

update-function:
	aws lambda update-function-code \
	  --function-name $(FUNCTION_NAME) \
	  --zip-file "fileb://$(ARTIFACT_ZIP)"
.PHONY: update-function

invoke-function:
	aws lambda invoke \
	  --function-name $(FUNCTION_NAME) \
	  --payload '{"text":"Hello"}' \
	  response.txt
.PHONY: invoke-function

update-and-check-function: update-function invoke-function
.PHONY: update-and-check-function
