.PHONY: build hugo clean dev

build: hugo
	go build -ldflags "-s -w" -o bin/jordi-codes .

hugo:
	hugo --config config/hugo.toml --minify

clean:
	rm -rf public bin

dev:
	hugo --config config/hugo.toml --minify
	go run .
