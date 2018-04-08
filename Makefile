run:
	go run ./cmd/fy/main.go $(t)

test:
	make run t=test

build:
	go build ./cmd/fy
