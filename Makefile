run-mintd:
	go run ./mintd

build-mintd:
	go build -o bin/mintd ./mintd

run-demo:
	go run ./demo

.PHONY: run-mintd build-mintd run-demo