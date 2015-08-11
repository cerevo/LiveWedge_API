package libvsw

import (
	"encoding/binary"
	"net"
	"unsafe"
)

type videoTransition struct {
	cmd          uint32
	cmdId        uint32
	param        uint32
	mode         uint8
	padding      [3]uint8
	main_src     uint8
	main_effect  uint8
	main_dip_src uint8
	padding1     uint8
	sub_src      uint8
	sub_effect   uint8
	sub_dip_src  uint8
	padding2     uint8
}

const VALUE_1 = (1 << 16)
const (
	VC_MODE_MAIN = iota
	VC_MODE_SUB
	VC_MODE_US
)

const (
	TRANSITION_TYPE_NULL = iota
	TRANSITION_TYPE_MIX
	TRANSITION_TYPE_DIP
	TRANSITION_TYPE_WIPE
	TRANSITION_TYPE_CUT = TRANSITION_TYPE_NULL
)
const (
	WIPE_HORIZONTAL   = iota
	WIPE_HORIZONTAL_R // _R means reversed pattern
	WIPE_VERTICAL
	WIPE_VERTICAL_R
	WIPE_HORIZONTAL_SLIDE
	WIPE_HORIZONTAL_SLIDE_R
	WIPE_VERTICAL_SLIDE
	WIPE_VERTICAL_SLIDE_R
	WIPE_HORIZONTAL_DOUBLE_SLIDE
	WIPE_HORIZONTAL_DOUBLE_SLIDE_R
	WIPE_VERTICAL_DOUBLE_SLIDE
	WIPE_VERTICAL_DOUBLE_SLIDE_R
	WIPE_SQUARE_TOP_LEFT /* top to bottom and left to right order */
	WIPE_SQUARE_TOP_LEFT_R
	WIPE_SQUARE_TOP
	WIPE_SQUARE_TOP_R
	WIPE_SQUARE_TOP_RIGHT
	WIPE_SQUARE_TOP_RIGHT_R
	WIPE_SQUARE_CENTER_LEFT
	WIPE_SQUARE_CENTER_LEFT_R
	WIPE_SQUARE_CENTER
	WIPE_SQUARE_CENTER_R
	WIPE_SQUARE_CENTER_RIGHT
	WIPE_SQUARE_CENTER_RIGHT_R
	WIPE_SQUARE_BOTTOM_LEFT
	WIPE_SQUARE_BOTTOM_LEFT_R
	WIPE_SQUARE_BOTTOM
	WIPE_SQUARE_BOTTOM_R
	WIPE_SQUARE_BOTTOM_RIGHT
	WIPE_SQUARE_BOTTOM_RIGHT_R
	WIPE_TYPE_NUM
)

func transMain(conn *net.TCPConn, param int, src int, effect int, dip int, manual int) {
	a := videoTransition{cmd: SW_ID_DoAutoSwitching,
		cmdId:        VALUE_1,
		param:        uint32(param),
		mode:         VC_MODE_MAIN,
		main_src:     uint8(src),
		main_effect:  uint8(effect),
		main_dip_src: uint8(dip)}
	//fmt.Printf("sizeof a=%d\n", unsafe.Sizeof(a))
	//buf := new(bytes.Buffer)
	//err := binary.Write(buf, LE, a)
	//checkError(err)
	//for _, b := range buf.Bytes() {
	//	fmt.Printf("%02x ", b)
	//}

	size := uint32(unsafe.Sizeof(a))
	err := binary.Write(conn, LE, size)
	checkError(err)
	err = binary.Write(conn, LE, a)
	checkError(err)
}

func transSub(conn *net.TCPConn, param int, src int, effect int, dip int, manual int) {
	a := videoTransition{cmd: SW_ID_DoAutoSwitching,
		cmdId:       VALUE_1,
		param:       uint32(param),
		mode:        VC_MODE_SUB,
		sub_src:     uint8(src),
		sub_effect:  uint8(effect),
		sub_dip_src: uint8(dip)}
	size := uint32(unsafe.Sizeof(a))
	err := binary.Write(conn, LE, size)
	checkError(err)
	err = binary.Write(conn, LE, a)
	checkError(err)
}

func transUs(conn *net.TCPConn, param int, src int, src2 int, effect int, dip int, manual int) {
	a := videoTransition{cmd: SW_ID_DoAutoSwitching,
		cmdId:        VALUE_1,
		param:        uint32(param),
		mode:         VC_MODE_US,
		main_src:     uint8(src),
		main_effect:  uint8(effect),
		main_dip_src: uint8(dip),
		sub_src:      uint8(src2)}
	size := uint32(unsafe.Sizeof(a))
	err := binary.Write(conn, LE, size)
	checkError(err)
	err = binary.Write(conn, LE, a)
	checkError(err)
}


// Cut changes the main screen to the specified src immediately.
func (vsw Vsw) Cut(src int) {
	//log.Printf("cut(%d)\n", src)
	if src < 0 || 4 < src {
		return
	}
	transMain(vsw.conn, 1, src, TRANSITION_TYPE_CUT, 0, 0)
}

// CutSub changes the sub screen to the specified src immediately.
func (vsw Vsw) CutSub(src int) {
	//log.Printf("cutSub(%d)\n", src)
	if src < 0 || 4 < src {
		return
	}
	transSub(vsw.conn, 1, src, TRANSITION_TYPE_CUT, 0, 0)
}

// CutUs changes both main and sub screen immediately.
func (vsw Vsw) CutUs(src int, src2 int) {
	//log.Printf("cutUs(%d,%d)\n", src, src2)
	if src < 0 || 4 < src {
		return
	}
	if src2 < 0 || 4 < src2 {
		return
	}
	if src == src2 {
		return
	}
	transUs(vsw.conn, 1, src, src2, TRANSITION_TYPE_CUT, 0, 0)
}

// Mix transits the main screen to the specified src.
//
// 
func (vsw Vsw) Mix(param int, src int) {
	//log.Printf("mix(%d, %d)\n", param, src)
	if src < 0 || 4 < src {
		return
	}
	transMain(vsw.conn, param, src, TRANSITION_TYPE_MIX, 0, 0)
}

// Dip transits the main screen to the specified src through dip_src in the specified duration.
func (vsw Vsw) Dip(param int, src int, dip_src int) {
	//log.Printf("dip(%d, %d, %d)\n", param, src, dip_src)
	if src < 0 || 4 < src {
		return
	}
	transMain(vsw.conn, param, src, TRANSITION_TYPE_DIP, dip_src, 0)
}

// Wipe transits the main screen to the specified src in the specified duration, using the specified wipe_type.
func (vsw Vsw) Wipe(param int, src int, wipe_type int) {
	//log.Printf("wipe(%d, %d, %d)\n", param, src, wipe_type)
	if src < 0 || 4 < src {
		return
	}
	if wipe_type < 0 || wipe_type >= WIPE_TYPE_NUM {
		return
	}
	transMain(vsw.conn, param, src, TRANSITION_TYPE_WIPE+wipe_type, 0, 0)
}
