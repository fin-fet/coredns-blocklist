# TODO: Improve Makefile to support other architectures
GO_FILES=$(shell find . -name '*.go')
ARCH=linux/arm64

default: $(ARCH)

linux/arm64: $(GO_FILES)
	mkdir -p build/linux/arm64
	GOOS=linux GOARCH=arm64 go build -o build/linux/arm64/coredns ./run/main.go

linux/amd64: $(GO_FILES)
	mkdir -p build/linux/amd64
	GOOS=linux GOARCH=amd64 go build -o build/linux/amd64/coredns ./run/main.go

clean:
	rm -rf ./build

run: $(ARCH)
	cp run/Corefile build/$(ARCH)/Corefile
	cp run/blocklist.txt build/$(ARCH)/blocklist.txt
	
	cd run && ../build/$(ARCH)/coredns -conf Corefile

dist: linux/arm64 linux/amd64