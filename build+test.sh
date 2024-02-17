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

if ! type -p shadow; then
  v go get     golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
  v go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
fi

echo date...
v go test -v -covermode=count -coverprofile=date.out .
v go tool cover -func=date.out

for d in clock period timespan view; do
  echo $d...
  v go test -v -covermode=count -coverprofile=$d.out ./$d
  v go tool cover -func=$d.out
done

v gofmt -l -w *.go */*.go

v go vet ./...

v shadow ./...

v go install ./datetool
