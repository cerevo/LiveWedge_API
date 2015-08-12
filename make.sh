#!/bin/sh
set -x

export GOPATH=$PWD:$GOPATH
(cd src/libvsw; go generate; go install)
(cd src/autotrans; go build)
(cd src/sample_trans; go build)
(cd src/sample_wipe; go build)
(cd src/sample_pinp; go build)
(cd src/sample_chromakey; go build)
(cd src/sample_rec; go build)
(cd src/sample_status; go build)
(cd src/find0; go build)
