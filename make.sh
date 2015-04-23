#!/bin/sh
set -x

export GOPATH=$PWD:$GOPATH
(cd src/libvsw; go install)
(cd src/autotrans; go build)
(cd src/rec0; go build)
(cd src/wipetest; go build)
