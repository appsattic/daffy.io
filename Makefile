all:
	echo 'Provide a target: server clean'

fmt:
	find src/ -name '*.go' -exec go fmt {} ';'

build: fmt
	gb build all

server: build
	./bin/server

test:
	gb test -v

clean:
	rm -rf bin/ pkg/

.PHONY: server
