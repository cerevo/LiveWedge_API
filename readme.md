# LiveWedge control library and samples

// Copyright 2015, Cerevo Inc. All rights reserved.  
// Use of this source code is governed by a BSD-style  
// license that can be found in the LICENSE file.  

This provides basic operations for the video switcher, "LiveWedge".

(See http://livewedge.cerevo.com/en/ about LiveWedge.)

This library is still alpha version. Compatiblity may break in future update.

## Supported operations
* Screen transfer: Cut, Mix, Dip, Wipe
* Sub screen control: PinP, Chroma-key
* Start and stop recording and broadcasting
* Upload a still picture and use it as ch.4 input source
* Find LiveWedge within the same network

## Not yet supported
Getting status from LiveWedge is still under construction. func (vsw Vsw) Request* are not yet fully documented.

## Contents
### libvsw
Common library for manupulating LiveWedge by network.

### autotrans
Sample program to make video transition automatically. See src/autotrans/00Readme.txt

### sample_trans
Very simple program just to repeat cut and mix.

### sample_wipe
Very simple program just to repeat wipe all pattern. 

### sample_pinp
Very simple program of Picture in Picture. 

### sample_chromakey
Very simple program of chroma key. 

### sample_rec
Very simple program just to send recording start/stop command with web UI.

### sample_status
A sample program for getting status via UDP.

### sample_find
A sample program for finding a LiveWedge within the same network.

## How to build

0. Install go language.
Tested in linux/amd64. Go version 1.4 and 1.5. I hope Mac/Windows works, too.

1. Execute ./make.sh at the top directory.

## API document

https://godoc.org/github.com/cerevo/LiveWedge_API/src/libvsw

