# Box Sync Client for Linux

## Development

First, install Go and set up your Go workspace by following [these instructions](https://golang.org/doc/code.html). In short:

```bash
mkdir ~/go
echo "export GOPATH=$HOME/go" >> ~/.bashrc
echo "export PATH=$GOPATH/bin:$PATH" >> ~/.bashrc
```

Next, clone this repo:

```bash
mkdir -p ~/go/src/gitlab-beta.engr.illinois.edu/sp-box
git clone git@gitlab-beta.engr.illinois.edu:sp-box/boxsync.git ~/go/src/gitlab-beta.engr.illinois.edu/sp-box/boxsync
```

To build:

```bash
cd ~/go/src/gitlab-beta.engr.illinois.edu/sp-box/boxsync
go install ./... # Build all main packages underneath the current directory and install in $GOPATH/bin
                 # Same as `go install gitlab-beta.engr.illinois.edu/sp-box/boxsync/cmd/boxsync` in this case
```

To run, just run `boxsync` because `$GOPATH/bin` is in your `$PATH`.

We will use [govendor](https://github.com/kardianos/govendor) for vendoring.

To add dependencies, a proper usage is like:

govendor fetch github.com/fsnotify/fsnotify

## Code style

In general, conform to the style guidelines described [here](https://github.com/golang/go/wiki/CodeReviewComments) and configure your text editor to run [gofmt](https://golang.org/cmd/gofmt/) on save. Project-specific code style guidelines are described below.

### Import order

Imports should be organized into the following three groups, with blank lines between them.

```go
import (
    // Standard library packages
    "testing"
    "time"

    // External (vendor) packages
    "github.com/stretchr/testify/assert"
    "golang.org/x/oauth2"

    // Internal packages
    "gitlab-beta.engr.illinois.edu/sp-box/boxsync/auth/mocks"
)
```

Imports in each group should be alphabetized; `gofmt` will do this automatically.

## Editors
### Vim plugins
I found following plugins can come handy for vim.

(1) Install Vundle for vim(https://github.com/VundleVim/Vundle.vim.git).

(2) Install YouCompleteMe for go (https://github.com/Valloric/YouCompleteMe.git).

(3) Install vim-go for go (https://github.com/fatih/vim-go.git).

(4) Install screen tool tmux.

(5) Install tagging tool gotags.

Don't forget to configure ~/.vimrc, various configurations can be found online.
