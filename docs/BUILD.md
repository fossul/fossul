# Build
The purpose is to provide guidelines for developers setting up dev environment. The instructions are for a Fedora 28 development environment but any Linux or MacOS should work.

## Install the Go programming language. 
```$ sudo dnf install -y go```

To build source code and setup a development ensure the following environment parameters are exported to the shell and set in user profile (.bashrc):
```
$ vi /home/ktenzer.bashrc
export GOPATH=/home/fedora/go
export GOBIN=/home/fedora
PATH=$PATH:$GOBIN
```

## Download dep binary. Dep is used for dependency and package management. Build scripts will call dep to download correct dependencies.
```$ curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh```

This will download and install dep into $GOBIN

## Clone the Fossul Github repository from '$GOPATH/src' in this case '/home/fedora/go/src'.
```$ git clone https://github.com/ktenzer/fossul.git```

## Change directory to the Fossul Github repository
```$ cd /home/fedora/go/src/fossul```

## Update Plugin Dir parameter in fossul build script
```
vi fossul-build.sh
PLUGIN_DIR="/home/fedora/plugins"
```

## Run the fossul build script
```$ /home/fedora/go/src/fossul/fossul-build.sh```