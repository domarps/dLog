# dLog
Distributed Log Querier
# dlog
Machine Programming I - Distributed Log Querier

## Setting Up Your Go Environment

This is, after all, a set of applications written in Go.
So first [install GVM](https://github.com/moovweb/gvm) - loan_surgeon requires
go1.4.1 or beyond to build correctly.
```
gvm install go1.4

gvm use go1.4
# if you use autoenv you may want to add the line above to your .env file
```

If you know what you are doing, feel free to clone `dlog` in the
correct go path location. If you are unsure, run:
```
go get gitlab-beta.engr.illinois.edu/mcconne7/dlog
```

Finally, run the following script to install all necessary
packages (including `golint`):

```
./go_get.sh
```

## Production

### Installation
Deploying to production:

```
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
source ~/.gvm/scripts/gvm
gvm install go1.4
gvm use go1.4
go get gitlab-beta.engr.illinois.edu/mcconne7/dlog
cd $GOPATH/src/gitlab-beta.engr.illinois.edu/mcconne7/dlog
./go_get.sh
go install ./...
```

## Development

### Dependencies
In order to build and run dlog you will need
* go
* ruby

### How to Build
To build dlog run the build script:
```
./build.rb
```
This:
1. formats the code
2. lints the code
3. tests the code
4. builds the code
5. installs the code

## Usage: Client
client: run

* You are probably going to want to end your commands with `2> /dev/null`
  to avoid log messages(to suppress the output). For example start the
  client with: `grep_client 127.0.0.1:3000 2>/dev/null`
* 127.0.0.1:3000 is the example URL

```
dgrep_client 127.0.0.1:3000
2015/09/12 20:28:07 client.go:46: Connection to: [127.0.0.1:3000]
```

### Status
returns the status of the server

```
status
{ "status": "running", "running": true }
```

### Echo
echos the command back to you

```
echo hello world!
hello world!
```

### System
runs the system command on all servers

```
sys grep go README.md
* go
go1.4.1 or beyond to build correctly.
gvm install go1.4
gvm use go1.4
correct go path location. If you are unsure, run:
go get gitlab-beta.engr.illinois.edu/mcconne7/dlog
packages (including `golint`):
./go_get.sh
```

### Log
runs a system command on each servers log file

```
log grep go
* go
go1.4.1 or beyond to build correctly.
gvm install go1.4
gvm use go1.4
correct go path location. If you are unsure, run:
go get gitlab-beta.engr.illinois.edu/mcconne7/dlog
packages (including `golint`):
./go_get.sh
```

## Usage: Server
server: run

### Start
```
dgrep_server 3000
dgrep server:2015/09/12 20:26:09 server.go:70: Server [127.0.0.1:3000]: is ending? false begins on: 127.0.0.1:3000
```
