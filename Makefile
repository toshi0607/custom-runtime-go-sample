build: build-function build-bootstrap

build-function:
	GOARCH=amd64 GOOS=linux go build -o artifacts/handler ./handlers

build-bootstrap:
	GOARCH=amd64 GOOS=linux go build -o artifacts/bootstrap ./init

zip:
	cp artifacts/handler handler
	cp artifacts/bootstrap bootstrap
	zip artifacts.zip handler bootstrap
	chmod +x handler bootstrap
	rm handler
	rm bootstrap

cp:
	cp bootstrap artifacts/bootstrap

deploy: build zip up

recf: d cf

cf: build zip
	aws lambda create-function \
	  --function-name "go-custome-runtime-sample" \
	  --zip-file "fileb://artifacts.zip" \
	  --handler "handler" \
	  --runtime provided \
	  --role $(LAMBDA_ROLL_ARN)

d:
	aws lambda delete-function \
	  --function-name "go-custome-runtime-sample"

up:
	aws lambda update-function-code \
	  --function-name "go-custome-runtime-sample" \
	  --zip-file "fileb://artifacts.zip"

li:
	aws lambda invoke \
	  --function-name "go-custome-runtime-sample" \
	  --payload '{"text":"Hello"}' \
	  response.txt

uc: up li
