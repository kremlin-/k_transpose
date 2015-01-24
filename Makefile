.PHONY: all

all: build

clean:
	go clean

install:
	install -s -o root -g root -m 555 k_transpose /usr/bin/k_transpose

build:
	go build 
