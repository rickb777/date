#!/bin/bash -e
# Developer tool to run the tests and obtain HTML coverage reports.

DIR=$PWD
DOT="$(dirname $0)"
cd $DOT
TOP=$PWD

# install Goveralls if absent
if ! type -p goveralls; then
  echo go install github.com/mattn/goveralls
  go install github.com/mattn/goveralls
fi

mkdir -p reports
rm -f reports/*.html coverage?*.*

for file in $(find . -type f -name \*_test.go | fgrep -v vendor/); do
  dirname $file >> coverage$$.tmp
done

sort coverage$$.tmp | uniq | tee coverage$$.dirs

for pkg in $(cat coverage$$.dirs); do
  name=$(echo $pkg | sed 's#^./##' | sed 's#/#-#g')
  [ "$pkg" = "." ] && name=$(basename $PWD)
  echo $pkg becomes $name
  go test -v -coverprofile coverage$$.data $pkg
  if [ -f coverage$$.data ]; then
    go tool cover -html coverage$$.data -o reports/$name.html
    unlink coverage$$.data
  fi
done

rm -f coverage$$.tmp coverage$$.dirs

if [ -n "$(type -p chromium-browser)" ]; then
  chromium-browser reports/*.html >/dev/null &
else
  ls -lh reports/
fi
