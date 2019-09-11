run:
	go run main.go

up:
	docker-compose -f local/docker-compose.yaml up

down:
	docker-compose -f local/docker-compose.yaml down

ps:
	docker-compose -f local/docker-compose.yaml ps

build:
	go build

release:
	goreleaser release --rm-dist --skip-publish --skip-sign
