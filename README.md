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
