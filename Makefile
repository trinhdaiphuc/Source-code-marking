build-binary:
	go build -o bin/server cmd/main.go

run-binary:
	./bin/server
