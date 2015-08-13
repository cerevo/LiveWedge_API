// Copyright 2015, Cerevo Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package libvsw

import (
	"encoding/binary"
	"io/ioutil"
	"unsafe"
)

type uploadFile0 struct {
	cmd           uint32
	contentLength uint32
	filename      [1024]byte
}

// UploadFile uploads the specified image file and use it as ch.4 input source.
//
// SD card has to be inserted to LiveWedge because it makes an intermediate file on SD card.
// The file has to be JPEG image file sized 1280 x 720.
// UploadFile returns error if the file can not read.
func (vsw Vsw) UploadFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
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
	return nil
}
