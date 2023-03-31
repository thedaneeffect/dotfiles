alias vi="hx"
alias zshrc="hx ~/.zshrc"
alias rc="hx ~/.rc"
alias todo="fixme *"

# Path to your oh-my-zsh installation.
export ZSH="$HOME/.oh-my-zsh"

export PATH="$PATH:$HOME/go/sdk/current/bin"
export PATH="$PATH:$HOME/go/bin"
export PATH="$PATH:$HOME/.local/bin"

export COLORTERM=truecolor
export GIT_EDITOR=hx
export GPG_TTY=$(tty)

ZSH_THEME="gruvbox"
SOLARIZED_THEME="dark"

plugins=(git)

source $ZSH/oh-my-zsh.sh
source $HOME/.rc

function tre() {
	tree -aC -I '.git|node_modules|.cache|.cargo' --dirsfirst "$@" | less -FRNX;
}

function mkfile() {
	mkdir -p `dirname "$1"`
	touch "$1"
}