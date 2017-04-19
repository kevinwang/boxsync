# Box Sync Client for Linux

[![build status](https://gitlab.engr.illinois.edu/sp-box/boxsync/badges/master/build.svg)](https://gitlab.engr.illinois.edu/sp-box/boxsync/commits/master)

## Development

First, install Go and set up your Go workspace by following [these instructions](https://golang.org/doc/code.html). In short:

```bash
mkdir ~/go
echo "export GOPATH=$HOME/go" >> ~/.bashrc
echo "export PATH=$GOPATH/bin:$PATH" >> ~/.bashrc
```

Next, clone this repo:

```bash
mkdir -p ~/go/src/gitlab.engr.illinois.edu/sp-box
git clone git@gitlab.engr.illinois.edu:sp-box/boxsync.git ~/go/src/gitlab.engr.illinois.edu/sp-box/boxsync
```

To build:

```bash
cd ~/go/src/gitlab.engr.illinois.edu/sp-box/boxsync
go install ./... # Build all main packages underneath the current directory and install in $GOPATH/bin
                 # Same as `go install gitlab.engr.illinois.edu/sp-box/boxsync/cmd/boxsync` in this case
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
    "gitlab.engr.illinois.edu/sp-box/boxsync/auth/mocks"
)
```

Imports in each group should be alphabetized; `gofmt` will do this automatically.

## Commit Message Format
You can put anything you want for commit message. But, in general, conform to a simple "subject-body" format, and use key action word as the first word for "subject". For example:

fixed/refactored/updated/removed/changed/released/merged/...  subsytem x for blab

subsystem x has a problem ..., and I did ... to fix...

more details...

## Editors
### Vim plugins
I found following plugins can come handy for vim.

(1) Install Vundle for vim(https://github.com/VundleVim/Vundle.vim.git).

(2) Install YouCompleteMe for go (https://github.com/Valloric/YouCompleteMe.git).

(3) Install vim-go for go (https://github.com/fatih/vim-go.git).

(4) Install screen tool tmux.

(5) Install tagging tool gotags.

Don't forget to configure ~/.vimrc, various configurations can be found online.

## Command line tool Manual

After building `boxcl`, command line tool is ready to use.

For more help information, `--help` is available for quick check.
```
$ boxcl [command_name] --help
```

Use command line example:
```
$ boxcl [command_name] [arguments...(will specify in following list)]
```

## Command line list

`user` - Allow user to login with OAuth for initialization. It will print the user id after login successes.

`dir` - Check for the availability of `$HOME/Box Sync` and make it the default sync folder. Print out the folder id.

`dlA` - Download all contents from cloud `Box Sync` folder.

`up [file_id]` - Upload file to Box root directory.

`up [file_id] [parent_folder_ id]` - Upload file to specific parent folder.

`upN [file_id] [file_local_src_path]` - Replace a specific file with new version.

`wE` -  Output event stream in real time.

`mkdir [folder_name]` - Create folder with `[folder_name]` in Box root directory.

`mkdir [folder_name] [parent_folder_ id]` - Create folder with `[folder_name]` in specific parent folder.

`rm [fild_id]` - Delete file.

`rmdir [folder_id]` - Delete folder recursively.

`ls` - List all files & folders in Box root directory.

`ls [parent_folder_id]` - List all files & folders in the parent folder.

