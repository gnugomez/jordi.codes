# jordi.codes

A personal site served over SSH, built with Go and the [Charmbracelet](https://charm.sh) stack.

## Visit

```sh
ssh -p 2222 jordi.codes
```

## Run locally

```sh
go run .
```

A `GITHUB_TOKEN` environment variable enables the GitHub contribution graph widget on the home screen.

## Deploy

```sh
docker compose up -d
```

The SSH host key is generated automatically on first start and persisted in the `ssh-host-key` volume so it survives container restarts.

## License

See [LICENSE](LICENSE).
