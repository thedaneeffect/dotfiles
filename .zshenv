# Sourced by EVERY zsh, including non-interactive and non-login shells.
# Keep this to PATH and environment only — no aliases, no output, no prompts.
#
# WHY this file exists: `ssh mbp 'tmux ...'` runs a NON-LOGIN, NON-INTERACTIVE
# shell. That skips both .zshrc and /etc/zprofile — and /etc/zprofile is where
# path_helper applies /etc/paths.d/homebrew. Without this file, Homebrew binaries
# (tmux included) are invisible to anything invoked over SSH.

typeset -U path PATH        # dedupe, since .zshrc prepends some of these again
export PATH="/opt/homebrew/bin:/opt/homebrew/sbin:$HOME/.local/bin:$PATH"
