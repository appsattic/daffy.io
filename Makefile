all:
	echo 'Provide a target: server clean'

fmt:
	find src/ -name '*.go' -exec go fmt {} ';'

vet:
	go vet src/internal/types/*.go
	go vet src/internal/store/*.go
	go vet src/cmd/server/*.go

staticcheck:
	GOPATH=/home/chilts/src/appsattic-daffy.io/vendor:/home/chilts/src/appsattic-daffy.io staticcheck src/internal/store/*.go
	GOPATH=/home/chilts/src/appsattic-daffy.io/vendor:/home/chilts/src/appsattic-daffy.io staticcheck src/internal/types/*.go
	GOPATH=/home/chilts/src/appsattic-daffy.io/vendor:/home/chilts/src/appsattic-daffy.io staticcheck src/cmd/server/*.go

test:
	gb test -v

build: fmt vet staticcheck test
	gb build all

server: build
	./bin/server

clean:
	rm -rf bin/ pkg/

.PHONY: server
