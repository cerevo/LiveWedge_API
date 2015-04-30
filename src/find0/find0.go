package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

const (
	SW_ID_FindSw    = 1
	SW_ID_FindSwAck = 2
)

const DISPLAY_NAME_LENGTH = 33

type AccessInfo struct {
	address     string
	udpPort     int
	tcpPort     int
	previewPort int
	displayName string
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

var LE = binary.LittleEndian

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
	//log.Println(len, "bytes read from", peer)
	buffer := bytes.NewBuffer(buf[:len])
	err = binary.Read(buffer, LE, &pkt)
	if err != nil {
		log.Fatal(err)
	}

	//log.Printf("adr=%s\n", peer.String())
	//log.Printf("cmd=%d\n", pkt.Cmd)
	//log.Printf("commandPort=%d\n", pkt.CommandPort)
	//log.Printf("tcpPort=%d\n", pkt.TcpPort)
	//log.Printf("previewPort=%d\n", pkt.PreviewPort)
	//log.Printf("isWiFiAP=%08x\n", pkt.IsWiFiAP)
	//log.Printf("displayName=%s\n", pkt.DisplayName)
	//log.Printf("\n")
	return pkt, peer, nil
}

func find0() []AccessInfo {
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
			udpPort:     int(pkt.CommandPort),
			tcpPort:     int(pkt.TcpPort),
			previewPort: int(pkt.PreviewPort),
		}
		a.address = strings.Split(peer.String(), ":")[0]
		n := bytes.Index(pkt.DisplayName[:], []byte{0})
		a.displayName = string(pkt.DisplayName[:n])
		pkts = append(pkts, a)
	}
	return pkts
}

func main() {
	items := find0()
	for _, v := range items {
		fmt.Printf("%#v\n", v)
	}
}
