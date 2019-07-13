BIN_DIR := $(GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter.exe
COCKROACH := ./db/init.local

localBuild: apiDocs lint database build run postman
	$(info localBuild complete)

lint: $(GOMETALINTER)
	$(info running linter)
	gometalinter ./...

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install &> /dev/null

database: $(COCKROACH)
	$(info setting up database)
	docker exec -it roach1 ./cockroach sql --insecure --execute="$(shell cat ./db/init_db.sql)"

$(COCKROACH):
	$(info setting up cockroach db)
	docker run -d --name=roach1 --hostname=roach1 --net=roachnet -p 26257:26257 -p 8080:8181 -v "${PWD}/db/cockroach-data" cockroachdb/cockroach:v2.1.0 start --insecure
	go get -u github.com/lib/pq
	touch ./db/init.local

tearDownDb:
	docker stop roach1
	docker rm roach1
	rm ./db/init.local
	unset SQL

build:
	$(info building)
	go build

dredd: run
	$(info running dredd)
	docker exec -it roach1 ./cockroach sql --insecure --execute="$(shell cat ./db/dredd_data_init.sql)"
	dredd docs/tkdo.apib http://localhost:7056/

postman: run
	$(info running postman)
	docker exec -it roach1 ./cockroach sql --insecure --execute="$(shell cat ./db/postman_data_init.sql)"
	newman run TKDO.postman_collection.json -e local-env.postman_environment.json

run: kill
	$(info running the server)
	nohup ./tkdo > nohup.out 2>&1 &

kill:
	$(info attempting to kill the server)
	if pgrep tkdo; then pkill tkdo; fi

apiDocs:
	$(info building API documentation)
	aglio -i ./docs/tkdo.apib -o ./docs/index.html

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
