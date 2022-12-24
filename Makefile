

build:
	go build -o bin/hosts

install: build
	cp bin/hosts ~/.bin/hosts