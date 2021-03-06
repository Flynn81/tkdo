BIN_DIR := $(GOPATH)/bin
COCKROACH := ./db/init.local

localBuild: apiDocs database lint unitTest build dredd postman godog run
	$(info localBuild complete)

lint:
	$(info running linter)
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.36.0
	golangci-lint run ./... --skip-files acceptance_test.go

running:
	ps -ef | grep tkdo

unitTest:
	go test -coverprofile=coverage.out ./...

coverage: unitTest
	go tool cover -html=coverage.out

godog:
ifndef TKDO_HOST
	$(error TKDO_HOST is not set)
endif
	docker exec -it tkdodb psql -U tk -d tkdo -c "$(shell cat ./db/clear_tables.sql)"
	go get github.com/cucumber/godog/cmd/godog
	godog
	docker exec -it tkdodb psql -U tk -d tkdo -c "$(shell cat ./db/clear_tables.sql)"

benchmarkAll: benchmarkHealthCheck benchmarkCreateUser benchmarkCreateTask benchmarkGetTasks
	$(info running all benchmarks)

benchmarkHealthCheck:
	$(info running health check benchmark)
	wrk -t 4 -c 10 -d 60 --latency --timeout 3s http://localhost:7056/hc

benchmarkCreateUser:
	$(info running create user benchmark)
	docker exec -it tkdodb psql -U tk -d tkdo -c "delete from task_user;"
	wrk -t 4 -c 10 -d 60 --latency --timeout 3s -s test/benchmark/create-user.lua http://localhost:7056/users
	docker exec -it tkdodb psql -U tk -d tkdo -c "delete from task_user;"

benchmarkCreateTask:
	$(info running create task benchmark)
	docker exec -it tkdodb psql -U tk -d tkdo -c "$(shell cat ./db/benchmark_create_task.sql)"
	wrk -t 4 -c 10 -d 60 --latency --timeout 3s -s test/benchmark/create-task.lua http://localhost:7056/tasks
	docker exec -it tkdodb psql -U tk -d tkdo -c "delete from task_user;"
	docker exec -it tkdodb psql -U tk -d tkdo -c "delete from task;"

benchmarkGetTasks:
	$(info running get tasks benchmark)
	docker exec -it tkdodb psql -U tk -d tkdo -c "$(shell cat ./db/benchmark_get_tasks.sql)"
	wrk -t 4 -c 10 -d 60 --latency --timeout 3s -s test/benchmark/get-tasks.lua http://localhost:7056/tasks
	docker exec -it tkdodb psql -U tk -d tkdo -c "delete from task_user;"
	docker exec -it tkdodb psql -U tk -d tkdo -c "delete from task;"

benchmark:
	$(info running benchmark)
	docker exec -it tkdodb psql -U tk -d tkdo -c "$(shell cat ./db/benchmark_get_tasks.sql)"
	wrk -t 4 -c 10 -d 60 --latency --timeout 3s -s test/benchmark/benchmark.lua http://localhost:7056
	docker exec -it tkdodb psql -U tk -d tkdo -c "delete from task_user;"
	docker exec -it tkdodb psql -U tk -d tkdo -c "delete from task;"

stress:
	go get -u github.com/tsenart/vegeta
	echo "GET http://localhost:7056/hc" | vegeta attack -duration=5s | tee results.bin | vegeta report

database: $(COCKROACH)
	$(info setting up database)

$(COCKROACH):
ifndef POSTGRES_PASSWORD
	$(error POSTGRES_PASSWORD is not set)
endif
ifndef DB_PASSWORD
		$(error DB_PASSWORD is not set)
endif
	$(info setting up db)
	docker pull postgres:9.6.20
	docker run --name tkdodb -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d -p 5432:5432 postgres:9.6.20
	sleep 10
	docker exec -it tkdodb psql -U postgres -c "CREATE ROLE tk LOGIN PASSWORD '${DB_PASSWORD}';"
	docker exec -it tkdodb psql -U postgres -c "$(shell cat ./db/init_db.sql)"
	docker exec -it tkdodb psql -U postgres -c "alter database tkdo owner to tk;"
	docker exec -it tkdodb psql -U tk -d tkdo -c "$(shell cat ./db/init_tables.sql)"
	go get -u github.com/lib/pq
	touch ./db/init.local

tearDownDb:
	docker stop tkdodb
	docker rm tkdodb
	rm ./db/init.local
	unset SQL

build:
	$(info building)
	go build

dredd: run
	$(info running dredd)
	docker exec -it tkdodb psql -U postgres -d tkdo -c "$(shell cat ./db/dredd_data_init.sql)"
	dredd docs/tkdo.apib http://localhost:7056/
	docker exec -it tkdodb psql -U postgres -d tkdo -c "DELETE FROM TASK; DELETE FROM TASK_USER;"

postman: run
	$(info running postman)
	docker exec -it tkdodb psql -U postgres -d tkdo -c "$(shell cat ./db/postman_data_init.sql)"
	newman run TKDO.postman_collection.json -e local-env.postman_environment.json
	docker exec -it tkdodb psql -U postgres -d tkdo -c "DELETE FROM TASK; DELETE FROM TASK_USER;"

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
	docker exec -it roach1 ./cockroach sql --insecure --execute="$(shell cat ./db/del_admin.sql)"

help:
	$(info targets are:)
	$(info localBuild)
	$(info lint)
	$(info build)
	$(info dredd)
	$(info postman)
	$(info run)
	$(info kill)
	$(info apiDocs)
