# ── build stage ──────────────────────────────────────────────────────────────
FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags "-s -w" -o /bin/jordi-codes .

# ── runtime stage ─────────────────────────────────────────────────────────────
FROM alpine:3.21

# openssh-keygen is needed only on first run to generate the host key;
# ca-certificates lets the GitHub API call succeed.
RUN apk add --no-cache openssh-keygen ca-certificates

WORKDIR /app
COPY --from=builder /bin/jordi-codes .

# The host key lives here. Mount a persistent volume to this path so the key
# survives container restarts and is shared across instances.
VOLUME ["/app/.ssh"]

# Generate a host key on first start if one doesn't already exist, then run.
ENTRYPOINT ["sh", "-c", "\
  mkdir -p /app/.ssh && \
  [ -f /app/.ssh/id_ed25519 ] || ssh-keygen -t ed25519 -f /app/.ssh/id_ed25519 -N '' -q && \
  exec /app/jordi-codes"]

EXPOSE 22
