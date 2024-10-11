BINARY_NAME=ynet

build:
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux ./cmd/y-net/main.go
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows ./cmd/y-net/main.go

clean:
	go clean
	rm ${BINARY_NAME}-linux
	rm ${BINARY_NAME}-windows

run:
	go run ./cmd/y-net/main.go

 test:
    go test -v -race -buildvcs ./...