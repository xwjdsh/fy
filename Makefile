run:
	go run ./cmd/fy/main.go $(t)

test:
	go test . ./bd ./by ./gg ./qq ./sg ./yd

build:
	go build ./cmd/fy

build-docker:
	docker build --build-arg VERSION=`git describe --tags` -t wendellsun/fy .

update-docker:
	make build-docker
	docker push wendellsun/fy
