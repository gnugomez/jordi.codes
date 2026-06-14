.PHONY: build hugo clean dev debug-hugo

build: hugo
	go build -ldflags "-s -w" -o bin/jordi-codes .

hugo:
	hugo --config config/hugo.toml --minify

debug-hugo:
	hugo --config config/hugo.toml

clean:
	rm -rf public bin

dev:
	hugo --config config/hugo.toml --minify
	go run .
