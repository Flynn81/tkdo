BIN_DIR := $(GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter.exe

localBuild: lint build run dredd postman
	$(info localBuild complete)

lint: $(GOMETALINTER)
	$(info running linter)
	gometalinter ./...

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install &> /dev/null

build:
	$(info building)
	go build

dredd: run
	$(info running dredd)
	dredd docs/tkdo.apib http://localhost:8080/

postman:
	$(info running postman)
	newman run TKDO.postman_collection.json

run: kill
	$(info running the server)
	nohup ./tkdo > /dev/null 2>&1 &

kill:
	$(info attempting to kill the server)
	if pgrep tkdo; then pkill tkdo; fi

help:
	$(info targets are:)
	$(info localBuild)
	$(info lint)
	$(info build)
	$(info dredd)
	$(info postman)
	$(info run)
	$(info kill)
