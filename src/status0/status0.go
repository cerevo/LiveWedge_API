package main

// Type 'go generate' to generate a.go.
//go:generate go run gen/gen.go

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"libvsw"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	SW_STATE_ID_StateMode                   = 100
	SW_STATE_ID_StateRecording              = 103
	SW_STATE_ID_StateFadeToDefaultColor     = 104
	SW_STATE_ID_StateExternalInput          = 105
	SW_STATE_ID_StateProgramOut             = 106
	SW_STATE_ID_StatePreviewOut             = 107
	SW_STATE_ID_StateDefaultBackgroundColor = 108
	SW_STATE_ID_StatePreviewMode            = 110
	SW_STATUS_ID_AudioMixer                 = 112
	SW_STATUS_ID_AudioMixerAll              = 113
	SW_STATUS_ID_AudioPeak                  = 114
	SW_STATUS_ID_VideoSwitcher              = 115
	SW_STATUS_ID_VideoSwitcherAuto          = 116
	SW_ID_MountNotify                       = 79
	SW_ID_CasterMessage                     = 85
	SW_ID_CasterStatistics                  = 86
)

// SwMode
const (
	SW_MODE_RTSP = iota
	SW_MODE_LIVE
	SW_MODE_RECORDING
)

// PreviewMode
const (
	PREVIEW_INPUT_TYPE_1 = iota
	PREVIEW_INPUT_TYPE_2
	PREVIEW_INPUT_TYPE_3
	PREVIEW_INPUT_TYPE_4
	PREVIEW_INPUT_TYPE_TILE
	PREVIEW_INPUT_TYPE_PROGRAM_OUT
	PREVIEW_INPUT_TYPE_SD = PREVIEW_INPUT_TYPE_4
)

// MountStatus
const (
	SW_MOUNT_UNMOUNTED = iota
	SW_MOUNT_READONLY
	SW_MOUNT_READWRITE
	SW_MOUNT_REQ_READONLY
	SW_MOUNT_REQ_READWRITE = SW_MOUNT_READWRITE
)

const (
	AUDIO_MIXER_CHANNEL_MASTER = iota
	AUDIO_MIXER_CHANNEL_INPUT1
	AUDIO_MIXER_CHANNEL_INPUT2
	AUDIO_MIXER_CHANNEL_INPUT3
	AUDIO_MIXER_CHANNEL_INPUT4
	AUDIO_MIXER_CHANNEL_CPU
	AUDIO_MIXER_CHANNEL_LINE
	AUDIO_MIXER_CHANNEL_NUM
)

const (
	AUDIO_MIXER_CATEGORY_GAIN = iota /* Master has no gain */
	AUDIO_MIXER_CATEGORY_VOLUME
	AUDIO_MIXER_CATEGORY_MUTE
	AUDIO_MIXER_CATEGORY_DELAY /* for only LINE */
)

type AudioMixerStatusType struct {
	Channel, Category uint8
	Value             uint16
}

type AudioMixerAllStatusType struct {
	Pairs       [AUDIO_MIXER_CHANNEL_NUM]struct{ Gain, Volume uint16 }
	Mute, Delay uint16
}

type AudioPeakStatusType struct {
	Peak [14]uint16
}

type RecordingStatusType struct {
	RecordingTime, RecordRemainTIme uint32
}

type FadeToDefaultColorStatusType struct {
	IsFade, AutoRemainTime uint32
}

type ExternalInputStatusType struct {
	InputType, PlayTiming, IsRepeat uint32
}

type ProgramOutStatusType struct {
	IsConnected   uint32
	Display       struct{ IsAuto, PixelHight, PixelWidth, Aspect, FrameRate uint32 }
	IsAudioEnable uint32
}

type PreviewOutStatusType struct {
	ProgramOutStatusType
	IsOSDEnable uint32
}

type Rect struct {
	X, Y int32
	W, H uint32
}

type SwitcherStatusType struct {
	Param, CmdId                         uint32
	Main_src, Sub_src, Sub_mode, Padding uint8
	Trans                                struct {
		Mode                                          uint8
		Padding                                       [3]uint8
		Main_src, Main_effect, Main_dip_src, Padding2 uint8
		Sub_src, Sub_effect, Sub_dip_src, Padding3    uint8
	}
	Chroma_floor, Chroma_ceil uint32
	Pinp                      struct {
		Scale Rect
		Crop  Rect
	}
	Pinp_border_color uint32
	Pinp_border_width uint8
}

type CasterMessageType struct {
	Category, Message uint8
	Stuff             [2]uint8
}

type CasterStatisticsType struct {
	Bitrate    uint32
	Queue      uint16
	Fps, Stuff uint8
}

var LE = binary.LittleEndian

var (
	SwMode                   uint32
	MountStatus              uint32
	PreviewMode              uint32
	DefaultBackgroundColor   uint32
	SwitcherStatus           SwitcherStatusType
	ProgramOutStatus         ProgramOutStatusType
	PreviewOutStatus         PreviewOutStatusType
	ExternalInputStatus      ExternalInputStatusType
	FadeToDefaultColorStatus FadeToDefaultColorStatusType
	RecordingStatus          RecordingStatusType
	AudioMixerAllStatus      AudioMixerAllStatusType
	AudioMixerStatus         AudioMixerStatusType
	AudioPeakStatus          AudioPeakStatusType
	CasterMessage            CasterMessageType
	CasterStatistics         CasterStatisticsType
)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "libvsw: Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func send(conn io.Writer, data []uint32) {
	err := binary.Write(conn, LE, data)
	checkError(err)
}

func openUdp(service string) *net.UDPConn {
	if strings.IndexRune(service, ':') < 0 {
		service += ":8888"
	}
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	checkError(err)

	conn, err := net.DialUDP("udp", nil, udpAddr)
	checkError(err)
	log.Println("connected")
	return conn
}

func read(conn io.Reader) {
	buf := make([]byte, 4096)
	len, err := conn.Read(buf)
	checkError(err)
	//log.Printf("len=%d\n", len)
	reader := bytes.NewReader(buf[:len])
	var cmd uint32
	err = binary.Read(reader, LE, &cmd)
	checkError(err)

	//log.Printf("cmd=%d\n", cmd)
	len -= 4
	switch cmd {
	case SW_STATE_ID_StateMode:
		readSwMode(len, reader)
	case SW_STATE_ID_StateRecording:
		readRecordingStatus(len, reader)
	case SW_STATE_ID_StateFadeToDefaultColor:
		readFadeToDefaultColorStatus(len, reader)
	case SW_STATE_ID_StateExternalInput:
		readExternalInputStatus(len, reader)
	case SW_STATE_ID_StateProgramOut:
		readProgramOutStatus(len, reader)
	case SW_STATE_ID_StatePreviewOut:
		readPreviewOutStatus(len, reader)
	case SW_STATE_ID_StateDefaultBackgroundColor:
		readDefaultBackgroundColor(len, reader)
	case SW_STATE_ID_StatePreviewMode:
		readPreviewMode(len, reader)
	case SW_STATUS_ID_AudioMixer:
		readAudioMixerStatus(len, reader)
	case SW_STATUS_ID_AudioMixerAll:
		readAudioMixerAllStatus(len, reader)
	case SW_STATUS_ID_AudioPeak:
		readAudioPeakStatus(len, reader)
	case SW_STATUS_ID_VideoSwitcher, SW_STATUS_ID_VideoSwitcherAuto:
		readSwitcherStatus(len, reader)
	case SW_ID_MountNotify:
		readMountStatus(len, reader)
	case SW_ID_CasterMessage:
		readCasterMessage(len, reader)
	case SW_ID_CasterStatistics:
		readCasterStatistics(len, reader)
	default:
		log.Printf("cmd=%d len=%d\n", cmd, len+4)
		for len > 0 {
			err = binary.Read(reader, LE, &cmd)
			checkError(err)
			log.Printf("%08x ", cmd)
			len -= 4
		}
		log.Printf("\n")
	}
}

func monitor(service string) {
	conn := openUdp(service)
	cmd := []uint32{33}
	send(conn, cmd)
	for {
		read(conn)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usge: %s IP_address_of_livewedge\n", os.Args[0])
		os.Exit(1)
	}
	vsw := libvsw.NewVsw(os.Args[1])

	go monitor(os.Args[1])
	rate := 10
	for {
		vsw.HeartBeat()
		time.Sleep(time.Duration(rate) * time.Second)
	}
}
