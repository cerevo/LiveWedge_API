// Package libvsw provides basic operations for the video switcher, "LiveWedge".
//
// This is still alpha version. Compatiblity may break in future update.
//
// Currently supported operations:
//   Screen transfer: Cut, Mix, Wipe
//   Sub screen control: PinP
//   Start and stop recording and broadcasting
//   Upload a still picture and use it as ch.4 input source
//
// Getting status from LiveWedge is still under construction.
// func (vsw Vsw) Request* are not yet documented.
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
	SW_ID_SwBasicInfo              = 0x03
	SW_ID_UploadFile               = 0x0c
	SW_ID_DoAutoSwitching          = 0x1c
	SW_ID_ChangeLiveBroadcastState = 0x12
	SW_ID_ChangeRecordingState     = 0x14
	SW_ID_TCPHeartBeat             = 0x23
	SW_ID_RecordingResult          = 0x25
	SW_ID_SetPinpGeometry          = 0x3d
	SW_ID_SetPinpBorder            = 0x3e
	SW_ID_SetChromaRange           = 0x3f
	SW_ID_SetSubMode               = 0x40
	SW_ID_SetTimezone              = 0x48
	SW_ID_SetTime                  = 0x49
	SW_ID_SetTimeAndZone           = 0x4a
	SW_ID_GetTimeAndZone           = 0x4b
)

// Vsw holds internal connection state.
type Vsw struct {
	conn    *net.TCPConn
	udpConn *net.UDPConn
	rev     int32
	update  int32
	mac     [8]uint8
}

var _vsw *Vsw

// Get firmware revision of LiveWedge
func (vsw Vsw) FirmwareRevision() int32 {
	return vsw.rev
}

// Get MAC address of LiveWedge
func (vsw Vsw) MacAddress() [8]uint8 {
	return vsw.mac
}

// Send heart beat to keep connection
func (vsw Vsw) HeartBeat() {
	buf := []uint32{SW_ID_TCPHeartBeat}
	send(vsw.conn, buf)

	var len int32
	err := binary.Read(vsw.conn, LE, &len)
	checkError(err)
	//fmt.Printf("len=%d\n", len)
	if len != 4 {
		fmt.Fprintf(os.Stderr, "libvsw: Fatal error. len != 4\n")
	}
	var cmd uint32
	err = binary.Read(vsw.conn, LE, &cmd)
	checkError(err)
	//fmt.Printf("cmd=%x\n", cmd)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "libvsw: Fatal error: %s\n", err.Error())
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

func readBasicInfo(vsw *Vsw) {
	var len int32
	err := binary.Read(vsw.conn, LE, &len)
	checkError(err)

	var cmd uint32
	err = binary.Read(vsw.conn, LE, &cmd)
	checkError(err)

	if cmd != SW_ID_SwBasicInfo {
		return
	}
	err = binary.Read(vsw.conn, LE, &vsw.rev)
	checkError(err)
	err = binary.Read(vsw.conn, LE, &vsw.update)
	checkError(err)
	err = binary.Read(vsw.conn, LE, &vsw.mac)
	checkError(err)
}

func openTcp(service string) (*net.TCPConn, error) {
	if strings.IndexRune(service, ':') < 0 {
		service += ":8888"
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	log.Println("TCP connected")
	return conn, nil
}

// NewVsw creates a new Vsw instance
//
// service: ip address or hostname of LiveWedge
// if failed to open LiveWedge, returns error.
func NewVsw(service string) (*Vsw, error) {
	if _vsw == nil {
		log.Println("New Vsw for", service)
		conn, err := openTcp(service)
		if err != nil {
			return nil, err
		}
		_vsw = new(Vsw)
		_vsw.conn = conn
		readBasicInfo(_vsw)
		_vsw.udpConn = openUdp(service)
		go monitorStatus(*_vsw)
	}
	return _vsw, nil
}

// Close the Vsw instance.
//
// Close TCP and UDP internal connections.
func (vsw Vsw) Close() {
	if _vsw != nil {
		_vsw.conn.Close()
		_vsw.udpConn.Close()
		_vsw = nil
	}
}
