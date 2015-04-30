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
