package libvsw

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"strings"
	"time"
)

// const (
// 	SW_ID_FindSw    = 1
// 	SW_ID_FindSwAck = 2
// )

const DISPLAY_NAME_LENGTH = 33

// AccessInfo is access information to LiveWedge.
// libvsw.Find() returns this.
type AccessInfo struct {
	Address     string
	UdpPort     int
	TcpPort     int
	PreviewPort int
	DisplayName string
}

type SwApi_FindSwAck struct {
	Cmd         uint32
	CommandPort uint16
	TcpPort     uint16
	PreviewPort uint16
	Dummy       uint16
	IsWiFiAP    uint32
	DisplayName [DISPLAY_NAME_LENGTH]byte
}

func read(conn net.PacketConn) (SwApi_FindSwAck, net.Addr, error) {
	var pkt SwApi_FindSwAck
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	buf := make([]byte, 4096)
	len, peer, err := conn.ReadFrom(buf)
	if err != nil {
		if opErr, ok := err.(*net.OpError); ok {
			if opErr.Timeout() {
				return pkt, peer, err
			}
		}
		log.Fatal(err)
	}
	buffer := bytes.NewBuffer(buf[:len])
	err = binary.Read(buffer, _LE, &pkt)
	if err != nil {
		log.Fatal(err)
	}
	return pkt, peer, nil
}

// Find LiveWedges in the same network.
func Find() []AccessInfo {
	conn, err := net.ListenPacket("udp", ":0")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	dst, err := net.ResolveUDPAddr("udp", "224.0.0.250:8888")
	if err != nil {
		log.Fatal(err)
	}

	// Send SW_ID_FindSw in little endian
	if _, err := conn.WriteTo([]byte{1, 0, 0, 0}, dst); err != nil {
		log.Fatal(err)
	}
	pkts := make([]AccessInfo, 0, 5)
	for {
		pkt, peer, err := read(conn)
		if err != nil {
			break
		}
		a := AccessInfo{
			UdpPort:     int(pkt.CommandPort),
			TcpPort:     int(pkt.TcpPort),
			PreviewPort: int(pkt.PreviewPort),
		}
		a.Address = strings.Split(peer.String(), ":")[0]
		n := bytes.Index(pkt.DisplayName[:], []byte{0})
		a.DisplayName = string(pkt.DisplayName[:n])
		pkts = append(pkts, a)
	}
	return pkts
}
