#!/usr/bin/env bash
# Spin up Ubuntu container for testing dotfiles setup

# Pass GITHUB_TOKEN if set (increases API rate limit for mise/homebrew)
DOCKER_ENV_ARGS=""
if [[ -n "${GITHUB_TOKEN}" ]]; then
  DOCKER_ENV_ARGS="-e GITHUB_TOKEN=${GITHUB_TOKEN}"
fi

docker run -it --rm \
  -v "$PWD:/dotfiles" \
  --name dotfiles-test \
  ${DOCKER_ENV_ARGS} \
  ubuntu:latest \
  bash -c '
    apt-get update -qq && apt-get install -y -qq curl git sudo build-essential procps file

    # Create non-root user (Homebrew requirement)
    useradd -m -s /bin/bash -G sudo testuser
    echo "testuser ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

    # Pre-create Homebrew directory with proper permissions
    mkdir -p /home/linuxbrew/.linuxbrew
    chown -R testuser:testuser /home/linuxbrew

    # Pass GITHUB_TOKEN to testuser environment
    if [[ -n "${GITHUB_TOKEN}" ]]; then
      echo "export GITHUB_TOKEN=${GITHUB_TOKEN}" >> /home/testuser/.bashrc
    fi

    echo ""
    echo "=== Dotfiles Test Environment ==="
    echo "User: testuser (with sudo)"
    echo "Directory: /dotfiles"
    echo "Homebrew path: /home/linuxbrew/.linuxbrew"
    if [[ -n "${GITHUB_TOKEN}" ]]; then
      echo "GitHub token: configured (increased rate limits)"
    fi
    echo ""
    echo "Run: ./setup.sh"
    echo ""

    # Switch to testuser in /dotfiles directory
    exec su - testuser -c "cd /dotfiles && exec bash"
  '
