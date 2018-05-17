test:
	go test . ./bd ./by ./gg ./qq ./sg ./yd

build:
	go build ./cmd/fy

run:
	go run ./cmd/fy/main.go $(t)

build-docker:
	docker build --build-arg VERSION=`git describe --tags` -t wendellsun/fy .

update-docker:
	make build-docker
	docker push wendellsun/fy
