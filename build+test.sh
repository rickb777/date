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

echo date...
v go test -v -covermode=count -coverprofile=date.out .
v go tool cover -func=date.out

for d in clock timespan view; do
  echo $d...
  v go test -v -covermode=count -coverprofile=$d.out ./$d
  v go tool cover -func=$d.out
done

v gofmt -l -w *.go */*.go

v go vet ./...

v go install ./datetool
