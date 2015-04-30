package libvsw

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

func readRecordingState(conn io.Reader) (state, result uint32) {
	var len int32
	err := binary.Read(conn, LE, &len)
	checkError(err)
	//fmt.Printf("len=%d\n", len)
	if len != 12 {
		fmt.Fprintf(os.Stderr, "libvsw: Fatal error \n")
	}

	var cmd uint32
	err = binary.Read(conn, LE, &cmd)
	checkError(err)
	//fmt.Printf("cmd=%x\n", cmd)
	if cmd != SW_ID_RecordingState {
		fmt.Fprintf(os.Stderr, "libvsw: Fatal error \n")
	}
	err = binary.Read(conn, LE, &state)
	checkError(err)
	err = binary.Read(conn, LE, &result)
	checkError(err)
	return state, result
}

func (vsw Vsw) ChangeLiveBroadcastState(mode int) {
	if mode != 0 && mode != 1 {
		return
	}
	sendKeyValue(vsw.conn, SW_ID_ChangeLiveBroadcastState, mode)
}

func (vsw Vsw) ChangeRecordingState(mode int) {
	if mode != 0 && mode != 1 {
		return
	}
	sendKeyValue(vsw.conn, SW_ID_ChangeRecordingState, mode)
	state, result := readRecordingState(vsw.conn)
	log.Printf("recording state=%d, result=%d\n", state, result)
}
