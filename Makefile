BIN_DIR := $(GOPATH)/bin
COCKROACH := ./db/init.local

.PHONY: help
help:
	@LC_ALL=C $(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'

buildbot: apiDocs lint unitTest build dredd postman godog run

everything: localBuild benchmarkAll stress
	$(info running everything)

localBuild: apiDocs tearDownDb database lint unitTest build dredd postman godog run
	$(info localBuild complete)

docker:
	$(info building docker image)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o tkdo
	docker build -t tkdo:latest -t tkdo:v0.0.2 -t ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/tkdo:latest -t ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/tkdo:v0.0.2 .

dockerRun:
	$(info running docker image)
	docker run -p ${TKDO_PORT}:${TKDO_PORT} --env TKDO_HOST --env TKDO_PORT --env TKDO_USER --env TKDO_PASSWORD --env TKDO_DBNAME --env TKDO_DYNAMOHOST --env TKDO_CORS tkdo:latest

dynamo:
	aws dynamodb list-tables --region us-east-2
	aws dynamodb create-table \
    --table-name user \
    --attribute-definitions \
				AttributeName=email,AttributeType=S \
    --key-schema \
        AttributeName=email,KeyType=HASH \
--provisioned-throughput \
        ReadCapacityUnits=10,WriteCapacityUnits=5
	aws dynamodb create-table \
    --table-name task \
    --attribute-definitions \
				AttributeName=id,AttributeType=S \
    --key-schema \
				AttributeName=id,KeyType=HASH \
--provisioned-throughput \
        ReadCapacityUnits=10,WriteCapacityUnits=5

ecrLogin:
	aws ecr get-login-password --region us-east-2 | docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com

ecrPush: ecrLogin
	docker push ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/tkdo:latest
	docker push ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/tkdo:v0.0.2

lint:
	$(info running linter)
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.36.0
	golangci-lint run ./... --skip-files acceptance_test.go --timeout 5m

running:
	ps -ef | grep tkdo

unitTest:
	go test -coverprofile=coverage.out ./...

coverage: unitTest
	go tool cover -html=coverage.out

zipForAws: zip buildforAWs
	zip -ur source.zip tkdo-for-aws

incrementVersion:
	$(eval foo := $(shell cat version.txt))
	@echo $(foo)
	$(eval v := $(shell ./version.sh $(foo) bug))
	@echo $(v)
	git tag $(v)
	@echo $(v) > version.txt

zip: incrementVersion
	rm source.zip
	zip -r source.zip */*.go ./*.go
	zip -ur source.zip version.txt

upload: zipForAws
	aws s3 cp ./source.zip s3://$(TKDO_S3)

godog:
ifndef TKDO_HOST
	$(error TKDO_HOST is not set)
endif
	go get github.com/cucumber/godog/cmd/godog@v0.12.0
	godog

benchmarkAll: benchmarkHealthCheck benchmarkCreateUser benchmarkCreateTask benchmarkGetTasks
	$(info running all benchmarks)

benchmarkHealthCheck:
	$(info running health check benchmark)
	wrk -t 4 -c 10 -d 60 --latency --timeout 3s http://localhost:7056/hc

benchmarkCreateUser:
	$(info running create user benchmark)
	wrk -t 4 -c 10 -d 60 --latency --timeout 3s -s test/benchmark/create-user.lua http://localhost:7056/users

benchmarkCreateTask:
	$(info running create task benchmark)
	wrk -t 4 -c 10 -d 60 --latency --timeout 3s -s test/benchmark/create-task.lua http://localhost:7056/tasks

benchmarkGetTasks:
	$(info running get tasks benchmark)
	wrk -t 4 -c 10 -d 60 --latency --timeout 3s -s test/benchmark/get-tasks.lua http://localhost:7056/tasks

benchmark:
	$(info running benchmark)
	wrk -t 4 -c 10 -d 60 --latency --timeout 3s -s test/benchmark/benchmark.lua http://localhost:7056

stress:
	go get -u github.com/tsenart/vegeta
	echo "GET http://localhost:7056/hc" | vegeta attack -duration=5s | tee results.bin | vegeta report

database: $(COCKROACH)
	$(info setting up database)
	cd db; docker-compose up -d;
	AWS_ACCESS_KEY_ID=X AWS_SECRET_ACCESS_KEY=X aws dynamodb list-tables --endpoint-url http://localhost:8000 --region x
	AWS_ACCESS_KEY_ID=X AWS_SECRET_ACCESS_KEY=X aws dynamodb create-table \
    --table-name user \
    --attribute-definitions \
				AttributeName=email,AttributeType=S \
    --key-schema \
        AttributeName=email,KeyType=HASH \
--provisioned-throughput \
        ReadCapacityUnits=10,WriteCapacityUnits=5 --region x --endpoint-url http://localhost:8000
	AWS_ACCESS_KEY_ID=X AWS_SECRET_ACCESS_KEY=X aws dynamodb create-table \
    --table-name task \
    --attribute-definitions \
				AttributeName=id,AttributeType=S \
    --key-schema \
				AttributeName=id,KeyType=HASH \
--provisioned-throughput \
        ReadCapacityUnits=10,WriteCapacityUnits=5 --region x --endpoint-url http://localhost:8000

$(COCKROACH):
	$(info setting up db)
	touch ./db/init.local

tearDownDb:
	rm ./db/init.local || true
	unset SQL
	AWS_ACCESS_KEY_ID=X AWS_SECRET_ACCESS_KEY=X aws dynamodb delete-table \
    --table-name user --region x --endpoint-url http://localhost:8000 || true
	AWS_ACCESS_KEY_ID=X AWS_SECRET_ACCESS_KEY=X aws dynamodb delete-table \
    --table-name task --region x --endpoint-url http://localhost:8000 || true
	cd db; docker-compose down;

buildForAws:
	$(info building for AWS)
	GOOS=linux go build -o tkdo-for-aws

build:
	$(info building)
	go build

dredd: run
	$(info running dredd)
	AWS_ACCESS_KEY_ID=X AWS_SECRET_ACCESS_KEY=X aws dynamodb put-item \
    --table-name user \
    --item '{"id": {"S":"00000000-0000-0000-0000-000000000000"},"name": {"S":"Pat Smith"},"email": {"S":"somethingelse@something.com"},"status": {"S":"status"}}'\
		--return-consumed-capacity TOTAL --region x --endpoint-url http://localhost:8000
	AWS_ACCESS_KEY_ID=X AWS_SECRET_ACCESS_KEY=X aws dynamodb put-item \
	    --table-name task \
	    --item '{"id": {"S":"60853a85-681d-4620-9677-946bbfdc8fbc"},"name": {"S":"clean the gutters"},"taskType": {"S":"basic|recurring"},"status": {"S":"new"},"userId":{"S":"00000000-0000-0000-0000-000000000000"}}'\
			--return-consumed-capacity TOTAL --region x --endpoint-url http://localhost:8000
	dredd docs/tkdo.apib http://localhost:7056/

postman: run
	$(info running postman)
	newman run TKDO.postman_collection.json -e local-env.postman_environment.json

run: kill
	$(info running the server)
	nohup ./tkdo > nohup.out 2>&1 &

kill:
	$(info attempting to kill the server)
	if pgrep tkdo; then pkill tkdo; fi

reload: kill build run
	$(info reloading local server -> kill -> build -> run)

apiDocs:
	$(info building API documentation)
	aglio -i ./docs/tkdo.apib -o ./docs/index.html

#this needs to be tested
delAdmin:
	$(info deleting admin users)
