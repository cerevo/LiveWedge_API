// Copyright 2015, Cerevo Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package libvsw

// Type 'go generate' to generate a.go.
//go:generate go run gen/gen.go

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
	"strings"
)

// ID for each status
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
	SW_STATUS_ID_SetPinpGeometry            = 117
	SW_STATUS_ID_SetPinpBorder              = 118
	SW_STATUS_ID_SetChromaRange             = 119
	SW_STATUS_ID_SetSubMode                 = 120

	SW_ID_MountNotify      = 79
	SW_ID_CasterMessage    = 85
	SW_ID_CasterStatistics = 86
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

// Audo mixer input
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

// Audio mixer category
const (
	AUDIO_MIXER_CATEGORY_GAIN = iota /* Master has no gain */
	AUDIO_MIXER_CATEGORY_VOLUME
	AUDIO_MIXER_CATEGORY_MUTE
	AUDIO_MIXER_CATEGORY_DELAY /* for only LINE */
)

type SwModeType struct {
	Cmd    uint32
	SwMode uint32
}
type MountStatusType struct {
	Cmd         uint32
	MountStatus uint32
}
type PreviewModeType struct {
	Cmd         uint32
	PreviewMode uint32
}
type DefaultBackgroundColorType struct {
	Cmd                    uint32
	DefaultBackgroundColor uint32
}

type AudioMixerStatusType struct {
	Cmd               uint32
	Channel, Category uint8
	Value             uint16
}

type AudioMixerAllStatusType struct {
	Cmd         uint32
	Pairs       [AUDIO_MIXER_CHANNEL_NUM]struct{ Gain, Volume uint16 }
	Mute, Delay uint16
}

type AudioPeakStatusType struct {
	Cmd  uint32
	Peak [14]uint16
}

type RecordingStatusType struct {
	Cmd                             uint32
	RecordingTime, RecordRemainTIme uint32
}

type FadeToDefaultColorStatusType struct {
	Cmd                    uint32
	IsFade, AutoRemainTime uint32
}

type ExternalInputStatusType struct {
	Cmd                             uint32
	InputType, PlayTiming, IsRepeat uint32
}

type ProgramOutStatusType struct {
	Cmd           uint32
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
	Cmd                                  uint32
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
	Cmd               uint32
	Category, Message uint8
	Stuff             [2]uint8
}

type CasterStatisticsType struct {
	Cmd        uint32
	Bitrate    uint32
	Queue      uint16
	Fps, Stuff uint8
}

type RecordingResultType struct {
	Cmd              uint32
	IsStateRecording uint32
	RecordingResult  uint32
}

type LiveBroadcastResultType struct {
	Cmd                 uint32
	DeviceId            [23]uint8
	Stuff               [1]uint8
	IsStateOnline       uint32
	LiveBroadcastResult uint32
}

type SubModeType struct {
	Cmd  uint32
	Mode uint32
}

type PinpGeometryType struct {
	Cmd     uint32
	Scale_x uint32
	Scale_y uint32
	Scale_w uint32
	Scale_h uint32
	Crop_x  uint32
	Crop_y  uint32
	Crop_w  uint32
	Crop_h  uint32
}

type PinpBorderType struct {
	Cmd   uint32
	Color uint32
	Width uint32
}

type ChromaRangeType struct {
	Cmd   uint32
	Floor uint32
	Ceil  uint32
}

func sendUdp(conn io.Writer, data []uint32) {
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
	log.Println("UDP connected")
	return conn
}

func readStatus(conn io.Reader) {
	buf := make([]byte, 4096)
	len, err := conn.Read(buf)
	checkError(err)
	//log.Printf("len=%d\n", len)
	reader := bytes.NewReader(buf[:len])
	reader2 := bytes.NewReader(buf[:4])
	var cmd uint32
	err = binary.Read(reader2, LE, &cmd)
	checkError(err)

	//log.Printf("cmd=%d\n", cmd)
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
	case SW_ID_RecordingResult:
		readRecordingResult(len, reader)
		RecordingResult.RecordingResult = 0
	case SW_STATUS_ID_SetSubMode:
		readSubMode(len, reader)
	case SW_STATUS_ID_SetPinpGeometry:
		readPinpGeometry(len, reader)
	case SW_STATUS_ID_SetPinpBorder:
		readPinpBorder(len, reader)
	case SW_STATUS_ID_SetChromaRange:
		readChromaRange(len, reader)
	default:
		log.Printf("readStatus: cmd=%d len=%d\n", cmd, len)
		for len > 0 {
			err = binary.Read(reader, LE, &cmd)
			checkError(err)
			log.Printf("%08x ", cmd)
			len -= 4
		}
		log.Printf("\n")
	}
}

func monitorStatus(vsw Vsw) {
	cmd := []uint32{33}
	sendUdp(vsw.udpConn, cmd)
	for {
		readStatus(vsw.udpConn)
	}
}
