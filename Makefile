build:
	go build

run: build
	./devtools

release:
	goreleaser release --rm-dist --skip-publish --skip-sign
