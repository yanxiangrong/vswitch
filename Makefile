export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct

all: fmt build

build: vswitchs vswitchc

fmt:
	go fmt ./...

vswitchs:
	env go build -o bin/vswitchs ./cmd/vswitchs

vswitchc:
	env go build -o bin/vswitchc ./cmd/vswitchc

clean:
	rm -f ./bin/vswitchs
	rm -f ./bin/vswitchc