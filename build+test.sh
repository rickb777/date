#!/bin/bash -e
cd "$(dirname $0)"
PATH=$HOME/go/bin:$PATH
unset GOPATH
export GOARCH=${1}

function v
{
  echo
  echo $@
  $@
}

if ! type -p goveralls; then
  v go install github.com/mattn/goveralls
fi

if ! type -p shadow; then
  v go get     golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
  v go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
fi

if ! type -p goreturns; then
  v go get     github.com/sqs/goreturns
  v go install github.com/sqs/goreturns
fi

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

v goreturns -l -w *.go */*.go

v go vet ./...

v shadow ./...

v go install ./datetool
