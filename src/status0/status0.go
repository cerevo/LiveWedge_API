package main

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
	"unsafe"
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
	DefautBackGroundColor    uint32
	SwitcherStatus           SwitcherStatusType
	ProgramOutStatus         ProgramOutStatusType
	PreviewOutStatus         PreviewOutStatusType
	ExternalInputStatus      ExternalInputStatusType
	FadeToDefaultColorStatus FadeToDefaultColorStatusType
	RecordingStatus          RecordingStatusType
	AudioMixerAllStatus      AudioMixerAllStatusType
	AudioMixerStatus         AudioMixerStatusType
	AudioPeakStatus          [14]uint16
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
		readAudioMixer(len, reader)
	case SW_STATUS_ID_AudioMixerAll:
		readAudioMixerAll(len, reader)
	case SW_STATUS_ID_AudioPeak:
		readAudioPeak(len, reader)
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

func readSwMode(len int, reader *bytes.Reader) {
	var a uint32
	if len != int(unsafe.Sizeof(a)) {
		log.Printf(" size mismatch %T len=%d\n", a, len)
		return
	}
	err := binary.Read(reader, LE, &a)
	checkError(err)
	if SwMode != a {
		log.Printf("Mode %#v\n", a)
		SwMode = a
	}

}

func readRecordingStatus(len int, reader *bytes.Reader) {
	var a RecordingStatusType
	if len != int(unsafe.Sizeof(a)) {
		log.Printf(" size mismatch %T len=%d\n", a, len)
		return
	}
	err := binary.Read(reader, LE, &a)
	checkError(err)
	if RecordingStatus != a {
		log.Printf("%#v\n", a)
		RecordingStatus = a
	}
}

func readFadeToDefaultColorStatus(len int, reader *bytes.Reader) {
	var a FadeToDefaultColorStatusType
	if len != int(unsafe.Sizeof(a)) {
		log.Printf(" size mismatch %T len=%d\n", a, len)
		return
	}
	err := binary.Read(reader, LE, &a)
	checkError(err)
	if FadeToDefaultColorStatus != a {
		log.Printf("%#v\n", a)
		FadeToDefaultColorStatus = a
	}
}

func readExternalInputStatus(len int, reader *bytes.Reader) {
	var a ExternalInputStatusType
	if len != int(unsafe.Sizeof(a)) {
		log.Printf(" size mismatch %T len=%d\n", a, len)
		return
	}
	err := binary.Read(reader, LE, &a)
	checkError(err)
	if ExternalInputStatus != a {
		log.Printf("%#v\n", a)
		ExternalInputStatus = a
	}
}

func readProgramOutStatus(len int, reader *bytes.Reader) {
	var a ProgramOutStatusType
	if len != int(unsafe.Sizeof(a)) {
		log.Printf(" size mismatch %T len=%d\n", a, len)
		return
	}
	err := binary.Read(reader, LE, &a)
	checkError(err)
	if ProgramOutStatus != a {
		log.Printf("%#v\n", a)
		ProgramOutStatus = a
	}
}

func readPreviewOutStatus(len int, reader *bytes.Reader) {
	var a PreviewOutStatusType
	if len != int(unsafe.Sizeof(a)) {
		log.Printf(" size mismatch %T len=%d\n", a, len)
		return
	}
	err := binary.Read(reader, LE, &a)
	checkError(err)
	if PreviewOutStatus != a {
		log.Printf("%#v\n", a)
		PreviewOutStatus = a
	}
}

func readDefaultBackgroundColor(len int, reader *bytes.Reader) {
	var a uint32
	if len != int(unsafe.Sizeof(a)) {
		log.Printf(" size mismatch %T len=%d\n", a, len)
		return
	}
	err := binary.Read(reader, LE, &a)
	checkError(err)
	if DefautBackGroundColor != a {
		log.Printf("DefaultBackgroundColor %#v\n", a)
		DefautBackGroundColor = a
	}
}

func readPreviewMode(len int, reader *bytes.Reader) {
	var a uint32
	if len != int(unsafe.Sizeof(a)) {
		log.Printf(" size mismatch %T len=%d\n", a, len)
		return
	}
	err := binary.Read(reader, LE, &a)
	checkError(err)
	if PreviewMode != a {
		log.Printf("PreviewMode %#v\n", a)
		PreviewMode = a
	}
}

func readMountStatus(len int, reader *bytes.Reader) {
	var a uint32
	if len != int(unsafe.Sizeof(a)) {
		log.Printf(" size mismatch %T len=%d\n", a, len)
		return
	}
	err := binary.Read(reader, LE, &a)
	checkError(err)
	if MountStatus != a {
		log.Printf("MountStatus %#v\n", a)
		MountStatus = a
	}
}

func readCasterMessage(len int, reader *bytes.Reader) {
	var a CasterMessageType
	if len != int(unsafe.Sizeof(a)) {
		log.Printf(" size mismatch %T len=%d\n", a, len)
		return
	}
	err := binary.Read(reader, LE, &a)
	checkError(err)
	if CasterMessage != a {
		log.Printf("%#v\n", a)
		CasterMessage = a
	}
}

func readCasterStatistics(len int, reader *bytes.Reader) {
	var a CasterStatisticsType
	if len != int(unsafe.Sizeof(a)) {
		log.Printf(" size mismatch %T len=%d\n", a, len)
		return
	}
	err := binary.Read(reader, LE, &a)
	checkError(err)
	if CasterStatistics != a {
		log.Printf("%#v\n", a)
		CasterStatistics = a
	}
}

func readSwitcherStatus(len int, reader *bytes.Reader) {
	var a SwitcherStatusType
	if len != int(unsafe.Sizeof(a)) {
		log.Printf("size mismatch %T len=%d\n", a, len)
		return
	}
	err := binary.Read(reader, LE, &a)
	checkError(err)
	if SwitcherStatus != a {
		log.Printf("%#v\n", a)
		SwitcherStatus = a
	}
}

func readAudioMixer(len int, reader *bytes.Reader) {
	var a AudioMixerStatusType
	if len != int(unsafe.Sizeof(a)) {
		return
	}
	err := binary.Read(reader, LE, &a)
	checkError(err)
	if AudioMixerStatus != a {
		log.Printf("%#v\n", a)
		AudioMixerStatus = a
	}
}

func readAudioMixerAll(len int, reader *bytes.Reader) {
	var a AudioMixerAllStatusType
	if len != int(unsafe.Sizeof(a)) {
		return
	}
	err := binary.Read(reader, LE, &a)
	checkError(err)
	if AudioMixerAllStatus != a {
		log.Printf("%#v\n", a)
		AudioMixerAllStatus = a
	}
}

func readAudioPeak(len int, reader *bytes.Reader) {
	var a [14]uint16
	if len != int(unsafe.Sizeof(a)) {
		log.Printf("size mismatch %d %#v\n", len, a)
		return
	}
	err := binary.Read(reader, LE, &a)
	checkError(err)
	if AudioPeakStatus != a {
		log.Printf("AudioPeak %#v\n", a)
		AudioPeakStatus = a
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
