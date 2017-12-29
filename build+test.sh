#!/bin/bash -e
cd $(dirname $0)
PATH=$HOME/gopath/bin:$GOPATH/bin:$PATH

if ! type -p goveralls; then
  echo go get github.com/mattn/goveralls
  go get github.com/mattn/goveralls
fi

echo date...
go test -v -covermode=count -coverprofile=date.out .
go tool cover -func=date.out
[ -z "$COVERALLS_TOKEN" ] || goveralls -coverprofile=date.out -service=travis-ci -repotoken $COVERALLS_TOKEN

for d in clock period timespan view; do
  echo $d...
  go test -v -covermode=count -coverprofile=$d.out ./$d
  go tool cover -func=$d.out
  [ -z "$COVERALLS_TOKEN" ] || goveralls -coverprofile=$d.out -service=travis-ci -repotoken $COVERALLS_TOKEN
done
