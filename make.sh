#!/bin/sh
set -x

if [ -z $GOPATH ]; then
    export GOPATH=$PWD
else
    export GOPATH=$PWD:$GOPATH
fi
(cd src/libvsw; go generate; go install)
(cd src/autotrans; go build)
(cd src/sample_trans; go build)
(cd src/sample_wipe; go build)
(cd src/sample_pinp; go build)
(cd src/sample_chromakey; go build)
(cd src/sample_rec; go build)
(cd src/sample_status; go build)
(cd src/sample_find; go build)
