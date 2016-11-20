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
