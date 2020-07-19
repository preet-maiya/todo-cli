.PHONY: build

all: build

build:
	go build -v -o todo main.go

install:
	go install -v

clean:
	git clean -fdx
