#!/bin/bash -e
cd "$(dirname $0)"
PATH=$HOME/go/bin:$PATH
unset GOPATH
export GO111MODULE=on
export GOARCH=${1}

function v
{
  echo
  echo $@
  $@
}

function require
{
  if [[ $1 = "-f" ]]; then
    rm -f ~/go/bin/$2
  fi
  if [[ ! -x ~/go/bin/$2 ]]; then
    v go install $4@$3
    cp -vf ~/go/bin/$2 ~/go/bin/$2.$3
  fi
  if [[ ! -x ~/go/bin/$2.$3 ]]; then
    v go install $4@$3
    cp -vf ~/go/bin/$2 ~/go/bin/$2.$3
  fi
}

require "$1" goimports  v0.1.0  golang.org/x/tools/cmd/goimports
require "$1" shadow     v0.1.0  golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
require "$1" goveralls  v0.0.7  github.com/mattn/goveralls

echo date...
v go test -v -covermode=count -coverprofile=date.out .
v go tool cover -func=date.out
#[ -z "$COVERALLS_TOKEN" ] || goveralls -coverprofile=date.out -service=travis-ci -repotoken $COVERALLS_TOKEN

for d in clock period timespan view; do
  echo $d...
  v go test -v -covermode=count -coverprofile=$d.out ./$d
  v go tool cover -func=$d.out
  #[ -z "$COVERALLS_TOKEN" ] || goveralls -coverprofile=$d.out -service=travis-ci -repotoken $COVERALLS_TOKEN
done

v goimports -l -w *.go */*.go

v go vet ./...

v shadow ./...

v go install ./datetool
