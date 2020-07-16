.PHONY: build

all: build

build:
	go build -v -o todo main.go

clean:
	git clean -fdx
