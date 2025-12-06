# ============================================================================
# XDG Base Directory Specification
# ============================================================================
export XDG_CONFIG_HOME="$HOME/.config"
export XDG_DATA_HOME="$HOME/.local/share"
export XDG_CACHE_HOME="$HOME/.cache"
export XDG_STATE_HOME="$HOME/.local/state"

# Respect XDG directories
export CARGO_HOME="$XDG_DATA_HOME/cargo"
export GNUPGHOME="$XDG_DATA_HOME/gnupg"
export LESSHISTFILE="$XDG_STATE_HOME/less/history"
export NPM_CONFIG_CACHE="$XDG_CACHE_HOME/npm"
export DOCKER_CONFIG="$XDG_CONFIG_HOME/docker"
export GRADLE_USER_HOME="$XDG_DATA_HOME/gradle"
export RUSTUP_HOME="$XDG_DATA_HOME/rustup"
export SDKMAN_DIR="$XDG_DATA_HOME/sdkman"
export USQL_HISTORY="$XDG_STATE_HOME/usql/history"

# ============================================================================
# Environment
# ============================================================================
export EDITOR=hx
export PATH="$HOME/.local/bin:$PATH"
export PATH="$HOME/.local/share/cargo/bin:$PATH"
export PATH="$HOME/.local/share/mise/shims:$PATH"
export PATH="$HOME/go/bin:$PATH"

# opencode
export PATH="$XDG_DATA_HOME/opencode/bin:$PATH"
export OPENCODE_DISABLE_AUTOUPDATE=true
export OPENCODE_DISABLE_AUTOCOMPACT=true
export UV_MANAGED_PYTHON=1

# History
export HISTFILE=~/.zsh_history
export HISTSIZE=10000
export SAVEHIST=20000
setopt HIST_IGNORE_ALL_DUPS
setopt HIST_IGNORE_DUPS
setopt SHARE_HISTORY

# Directory navigation
setopt AUTO_PUSHD           # Automatically push directories onto stack
setopt CDABLE_VARS          # Allow cd to variable names
setopt PUSHD_MINUS          # Swap meaning of cd +1 and cd -1
setopt PUSHD_IGNORE_DUPS    # Don't push duplicate directories
setopt PUSHD_SILENT         # Don't print directory stack after pushd/popd

# Persistent directory stack across sessions
DIRSTACKFILE="$XDG_STATE_HOME/zsh/dirs"
if [[ -f "$DIRSTACKFILE" ]] && [[ ${#dirstack[*]} -eq 0 ]]; then
    dirstack=( ${(f)"$(< $DIRSTACKFILE)"} )
    [[ -d "${dirstack[1]}" ]] && cd -q "${dirstack[1]}"
fi
chpwd() {
    print -l $PWD ${(u)dirstack} > "$DIRSTACKFILE"
}

DIRSTACKSIZE=20

# ============================================================================
# Homebrew
# ============================================================================
if [[ -f /opt/homebrew/bin/brew ]]; then
    eval "$(/opt/homebrew/bin/brew shellenv)"
elif [[ -f /usr/local/bin/brew ]]; then
    eval "$(/usr/local/bin/brew shellenv)"
elif [[ -f /home/linuxbrew/.linuxbrew/bin/brew ]]; then
    eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"
fi
export HOMEBREW_NO_ENV_HINTS=1
export HOMEBREW_NO_ANALYTICS=1

# ============================================================================
# mise - unified tool version management
# ============================================================================
eval "$(mise activate zsh 2>/dev/null)"

# ============================================================================
# Bun
# ============================================================================
export BUN_INSTALL="$XDG_DATA_HOME/bun"
export PATH="$BUN_INSTALL/bin:$PATH"
[ -s "$BUN_INSTALL/_bun" ] && source "$BUN_INSTALL/_bun"

# ============================================================================
# Completion & Interactive Features
# ============================================================================
export ZSH_COMPDUMP="$XDG_CACHE_HOME/zsh/zcompdump"
autoload -Uz compinit
compinit -d "$ZSH_COMPDUMP"

zstyle ':completion:*:directory-stack' list-colors '=(#b) #([0-9]#)*( *)==95=38;5;12'

# Key bindings (Home/End)
bindkey "^[[H" beginning-of-line
bindkey "^[[F" end-of-line

# Prompt (starship)
eval "$(starship init zsh)"

# Tool completions
eval "$(fzf --zsh)"
eval "$(zoxide init zsh)"
eval "$(mise completion zsh)"
eval "$(golangci-lint completion zsh)"

# ============================================================================
# Aliases
# ============================================================================

# Navigation
alias ..='cd ..'
alias ...='cd ../..'
alias .rc='. ~/.zshrc'

# Modern replacements
alias l='eza -la'
alias ls='eza -la'
alias ll='eza -la'
alias la='eza -la'
alias tree='eza --tree'

# Git shortcuts
alias gst='git status'
alias ga='git add'
alias gc='git commit'
alias gd='git diff'
alias gds='git diff --staged'
alias gl='git log'
alias gcob='git for-each-ref --format="%(refname:lstrip=2)" refs/heads/ refs/remotes/ | sed "s#^origin/##" | sort -u | fzf | xargs git switch'
alias glf='git log --oneline | fzf --preview "git show {1}"'
alias gp='git push'

# mise shortcuts
alias mi='mise'
alias mii='mise install'
alias miu='mise upgrade'
alias mis='mise use'
alias mil='mise list'
alias mio='mise outdated'
alias mt='mise task'

# AI CLIs
alias pup='uvx --from /tmp/code_puppy pup'

# Utilities
alias grep='grep --color=auto'
alias rc='$EDITOR ~/.$(basename $SHELL)rc'
alias myip='curl -s ifconfig.me'
alias ports='lsof -i -P -n | grep LISTEN'
alias bootstrap='bash <(curl -fsSL https://coffer.medieval.software/bootstrap) && source ~/.zshrc'

# ============================================================================
# System-specific overrides
# ============================================================================
# Source additional config files from ~/.config/zsh.d/ for system-specific
# settings that override the defaults above (e.g., secrets, Docker, etc.)
if [[ -d "$HOME/.config/zsh.d" ]]; then
    for config in "$HOME/.config/zsh.d"/*.zsh(N); do
        source "$config"
    done
fi

