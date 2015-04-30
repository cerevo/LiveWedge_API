package libvsw

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

var LE = binary.LittleEndian

const (
	SW_ID_UploadFile               = 0x0c
	SW_ID_DoAutoSwitching          = 0x1c
	SW_ID_ChangeLiveBroadcastState = 0x12
	SW_ID_ChangeRecordingState     = 0x14
	SW_ID_RecordingState           = 0x25
	SW_ID_SetPinpGeometry          = 0x3d
	SW_ID_SetSubMode               = 0x40
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

func NewVsw(service string) Vsw {
	log.Println("New Vsw for", service)
	vsw := Vsw{}
	vsw.conn = openTcp(service)
	read(vsw.conn)
	return vsw
}
