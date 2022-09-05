GO_FILES=$(shell find . -name '*.go')
OS ?= linux
ARCH ?= amd64
PLATFORM=$(OS)/$(ARCH)

default: build

build: $(GO_FILES)
	mkdir -p build/$(PLATFORM)
	GOOS=$(OS) GOARCH=$(ARCH) go build -o build/$(PLATFORM)/coredns ./run/main.go

run: build
	cp run/Corefile build/$(PLATFORM)/Corefile
	cp run/blocklist.txt build/$(PLATFORM)/blocklist.txt
	
	cd run && ../build/$(PLATFORM)/coredns -conf Corefile

clean:
	rm -rf ./build