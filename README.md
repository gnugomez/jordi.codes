# jordi.codes

A personal site served over SSH, built with Go and the [Charmbracelet](https://charm.sh) stack.

## Visit

```sh
ssh jordi.codes
```

## Run locally

```sh
go run .
```

A `GITHUB_TOKEN` environment variable enables the GitHub contribution graph widget on the home screen.

## Deploy to the edge (Fly.io)

This project uses [Fly.io](https://fly.io) to run the container across its global edge network. Fly supports raw TCP on port 22 natively — no port remapping needed.

**Prerequisites:** [install `flyctl`](https://fly.io/docs/hands-on/install-flyctl/) and run `fly auth login`.

```sh
# 1. Create the app (only needed once)
fly apps create jordi-codes

# 2. Create a persistent volume for the SSH host key (only needed once)
#    The key is generated automatically on first start and shared across instances.
fly volumes create ssh_host_key --size 1 --region mad

# 3. Set the GitHub token secret (optional)
fly secrets set GITHUB_TOKEN=your_token_here

# 4. Deploy
fly deploy
```

The `fly.toml` in this repo is pre-configured. The `ssh_host_key` volume is mounted at `/app/.ssh` — the key is generated on first start and reused on every subsequent deploy and across all instances.

## License

See [LICENSE](LICENSE).
