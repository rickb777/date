#!/bin/bash -ex
cd "$(dirname "$0")"
go install tool
mage install coverage crosscompile
cat report.out
