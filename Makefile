build:
	go build

release:
	goreleaser release --rm-dist --skip-publish --skip-sign
