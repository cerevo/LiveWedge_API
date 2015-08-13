// Copyright 2015, Cerevo Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package libvsw

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"unsafe"
)

func readRecordingState(conn io.Reader) (state, result uint32) {
	var len int32
	err := binary.Read(conn, LE, &len)
	checkError(err)
	//fmt.Printf("len=%d\n", len)
	if len != 12 {
		fmt.Fprintf(os.Stderr, "libvsw: Fatal error. len != 12 \n")
	}

	var cmd uint32
	err = binary.Read(conn, LE, &cmd)
	checkError(err)
	//fmt.Printf("cmd=%x\n", cmd)
	if cmd != SW_ID_RecordingResult {
		fmt.Fprintf(os.Stderr, "libvsw: Fatal error. cmd != SW_ID_RecrodingResult \n")
	}
	err = binary.Read(conn, LE, &state)
	checkError(err)
	err = binary.Read(conn, LE, &result)
	checkError(err)
	return state, result
}

func readLiveBroadcastResult(conn io.Reader) {
	var a LiveBroadcastResultType
	var len int32
	err := binary.Read(conn, LE, &len)
	checkError(err)
	//fmt.Printf("len=%d\n", len)

	if len != int32(unsafe.Sizeof(a)) {
		log.Printf(" size mismatch %T len=%d\n", a, len)
		return
	}
	err = binary.Read(conn, LE, &a)
	checkError(err)
	//log.Printf("%#v\n", a)
}

// ChangeLiveBroadcastState
//
// mode: 0 stop broadcasting, 1 start broadcasting
func (vsw Vsw) ChangeLiveBroadcastState(mode int) {
	if mode != 0 && mode != 1 {
		return
	}
	sendKeyValue(vsw.conn, SW_ID_ChangeLiveBroadcastState, mode)
	readLiveBroadcastResult(vsw.conn)
}

// ChangeRecordingState
//
// mode: 0 stop recording, 1 start recording
func (vsw Vsw) ChangeRecordingState(mode int) {
	if mode != 0 && mode != 1 {
		return
	}
	sendKeyValue(vsw.conn, SW_ID_ChangeRecordingState, mode)
	state, result := readRecordingState(vsw.conn)
	log.Printf("recording state=%d, result=%d\n", state, result)
}
