BUILDPATH=./cmd/dhcpserver

build:
	go get $(BUILDPATH)
	go build $(BUILDPATH)

all: build
