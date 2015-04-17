# LiveWedge API sample

## Contents
### libvsw
Common library for manupulating LiveWedge by network.


### autotrans
Sample program to make video transition automatically. See src/autotrans/00Readme.txt

### rec0
Very simple sample program just to send recording start/stop command.

## How to build

0. Install go language.
Tested in linux/amd64. Go version 1.4. I hope Mac/Windows works, too.

1. Add GOPATH environment variable

	export GOPATH=$PWD:$GOPATH

2. Install common library

	(cd src/libvsw; go install)

3. Build sample programs

	(cd src/autotrans; go build)  
	(cd src/rec0; go build)
