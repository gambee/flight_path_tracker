deps := $(shell find . -name *.go)

build/main: $(deps)
	go build -o build/main cmd/main/main.go

run: build/main
	./build/main

image: build/main
	docker build -t flight_path_tracker .

run_image: image
	docker run -p 8080:8080 flight_path_tracker

test:
	go run ./cmd/test/test.go
