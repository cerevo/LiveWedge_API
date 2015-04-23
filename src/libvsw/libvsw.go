package libvsw

import (
	//"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	//"time"
	"unsafe"
)

var LE = binary.LittleEndian

const (
	SW_ID_UploadFile               = 12
	SW_ID_DoAutoSwitching          = 0x1c
	SW_ID_ChangeLiveBroadcastState = 0x12
	SW_ID_ChangeRecordingState     = 0x14
	SW_ID_RecordingState           = 0x25
)

//const SW_ID_SetTimezone = 0x48
//const SW_ID_SetTime = 0x49
//const SW_ID_SetTimeAndZone = 0x4a
//const SW_ID_GetTimeAndZone = 0x4b

type Vsw struct {
	conn *net.TCPConn
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "libvsw: Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func send(conn io.Writer, data []uint32) {
	size := uint32(len(data) * 4)
	err := binary.Write(conn, LE, size)
	checkError(err)
	err = binary.Write(conn, LE, data)
	checkError(err)
}

func sendKeyValue(conn *net.TCPConn, key uint32, val int) {
	buf := [2]uint32{key, uint32(val)}
	send(conn, buf[:])
}

func readBasicInfo(conn io.Reader) {
	//fmt.Printf("read BasicInfo: ")
	var (
		rev    int32
		update int32
		mac    [8]uint8
	)
	err := binary.Read(conn, LE, &rev)
	checkError(err)
	err = binary.Read(conn, LE, &update)
	checkError(err)
	err = binary.Read(conn, LE, &mac)
	checkError(err)
	//fmt.Printf("rev=%d update=%d ", rev, update)
	//fmt.Printf("mac=%02x:%02x:%02x:%02x:%02x:%02x\n", mac[5], mac[4], mac[3], mac[2], mac[1], mac[0])
}

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

func read(conn io.Reader) {
	var len int32
	err := binary.Read(conn, LE, &len)
	checkError(err)
	//fmt.Printf("len=%d\n", len)

	var cmd uint32
	err = binary.Read(conn, LE, &cmd)
	checkError(err)

	switch cmd {
	case 35:
		fmt.Printf("TCPHeartBeat\n")
	case 3:
		readBasicInfo(conn)
	//case 0x7c:
	//	readTimeAndZone(conn)
	default:
		fmt.Printf("cmd=%08x ", cmd)
		len -= 4
		for len > 0 {
			err = binary.Read(conn, LE, &cmd)
			checkError(err)
			fmt.Printf("%08x ", cmd)
			len -= 4
		}
		fmt.Printf("\n")
	}
}

type uploadFile0 struct {
	cmd           uint32
	contentLength uint32
	filename      [1024]byte
}

func (vsw Vsw) UploadFile(filename string) {
	data, err := ioutil.ReadFile(filename)
	checkError(err)
	//fmt.Printf("len(data)=%d\n", len(data))
	a := uploadFile0{cmd: SW_ID_UploadFile,
		contentLength: uint32(len(data))}
	for i, c := range filename {
		a.filename[i] = byte(c)
	}
	size := uint32(unsafe.Sizeof(a)) + uint32(len(data))
	err = binary.Write(vsw.conn, LE, size)
	checkError(err)
	err = binary.Write(vsw.conn, LE, a)
	checkError(err)
	err = binary.Write(vsw.conn, LE, data)
	checkError(err)
}

type videoTransition struct {
	cmd          uint32
	value        uint32
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
const VC_MODE_MAIN = 0

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
		value:        VALUE_1,
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

func openTcp(service string) *net.TCPConn {
	if strings.IndexRune(service, ':') < 0 {
		service += ":8888"
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	log.Println("connected")
	return conn
}

func (vsw Vsw) Cut(src int) {
	//log.Printf("cut(%d)\n", src)
	if src < 1 || 4 < src {
		return
	}
	transMain(vsw.conn, 1, src, TRANSITION_TYPE_CUT, 0, 0)
}

func (vsw Vsw) Mix(param int, src int) {
	//log.Printf("mix(%d, %d)\n", param, src)
	if src < 1 || 4 < src {
		return
	}
	transMain(vsw.conn, param, src, TRANSITION_TYPE_MIX, 0, 0)
}

func (vsw Vsw) Dip(param int, src int, dip_src int) {
	//log.Printf("dip(%d, %d, %d)\n", param, src, dip_src)
	if src < 1 || 4 < src {
		return
	}
	transMain(vsw.conn, param, src, TRANSITION_TYPE_DIP, dip_src, 0)
}

func (vsw Vsw) Wipe(param int, src int, wipe_type int) {
	//log.Printf("wipe(%d, %d, %d)\n", param, src, wipe_type)
	if src < 1 || 4 < src {
		return
	}
	if wipe_type < 0 || wipe_type >= WIPE_TYPE_NUM {
		return
	}
	transMain(vsw.conn, param, src, TRANSITION_TYPE_WIPE+wipe_type, 0, 0)
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

func NewVsw(service string) Vsw {
	log.Println("New Vsw for", service)
	vsw := Vsw{}
	vsw.conn = openTcp(service)
	read(vsw.conn)
	return vsw
}
