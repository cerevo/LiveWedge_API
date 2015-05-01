#!/bin/sh
set -x

export GOPATH=$PWD:$GOPATH
(cd src/libvsw; go generate; go install)
(cd src/autotrans; go build)
(cd src/rec0; go build)
(cd src/status0; go build)
(cd src/wipetest; go build)
(cd src/pinptest; go build)
(cd src/find0; go build)
