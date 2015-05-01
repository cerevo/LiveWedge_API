# LiveWedge API sample

## Contents
### libvsw
Common library for manupulating LiveWedge by network.


### autotrans
Sample program to make video transition automatically. See src/autotrans/00Readme.txt

### rec0
Very simple program just to send recording start/stop command.

### wipetest
Very simple program just to repeat wipe all pattern. No web UI.

## How to build

0. Install go language.
Tested in linux/amd64. Go version 1.4. I hope Mac/Windows works, too.

1. Add GOPATH environment variable

	export GOPATH=$PWD:$GOPATH

2. Install common library

	(cd src/libvsw; go generate; go install)

3. Build sample programs at each directry

	(cd src/autotrans; go build)  
	(cd src/rec0; go build)  
	(cd src/wipetest; go build)  
	...  


Or just execute ./make.sh
