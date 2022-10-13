export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct

all: fmt build


build: vswitchs vswitchc

cross-build: cross-linux vswitchs vswitchc

cross-linux:
	env CGO_ENABLED=0
	env GOOS=linux
	env GOARCH=amd64

fmt:
	go fmt ./...

vswitchs:
	go build -o out/vswitchs ./cmd/vswitchs

vswitchc:
	go build -o out/vswitchc ./cmd/vswitchc

clean:
	rm -f ./bin/vswitchs
	rm -f ./bin/vswitchc