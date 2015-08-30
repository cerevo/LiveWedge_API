// Copyright 2015, Cerevo Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package libvsw

// Predefined chroma key enumuration
const (
	CHROMA_KEY_GREEN = iota
	CHROMA_KEY_BLUE
	CHROMA_KEY_PURPLE
	CHROMA_KEY_RED

	CHROMA_KEY_NUM
)

// Color range of chroma key
type chroma_range struct {
	floor, ceil int
}

var chroma_range_table = []chroma_range{
	chroma_range_entry(0, 0, 8-1, 8-2),
	chroma_range_entry(8+1, 0, 15, 8-2),
	chroma_range_entry(8+1, 8+2, 15, 15),
	chroma_range_entry(0, 8+2, 8-1, 15),
}

func chroma_range_entry(u0, v0, u1, v1 int) chroma_range {
	return chroma_range{
		(u0)<<4 | 1<<12 | (v0)<<20,
		(u1)<<4 | 14<<12 | (v1)<<20,
	}
}

// SetChromaRange sets color range of chroma key
func (vsw Vsw) SetChromaRange(floor, ceil int) {
	cmd := []uint32{SW_ID_SetChromaRange, uint32(floor), uint32(ceil)}
	send(vsw.conn, cmd)
}

// SetChromaKey sets chroma key in predefined color
// color is libvsw.CHROMA_KEY_GREEN, libvsw.CHROMA_KEY_BLUE, libvsw.CHROMA_KEY_PURPLE or libvsw.CHROMA_KEY_RED.
// For other color, use SetChromaRange and then call SetSubMode(SUB_MODE_CHROMAKEY)
func (vsw Vsw) SetChromaKey(color int) {
	if 0 > color || color >= CHROMA_KEY_NUM {
		return
	}
	vsw.SetChromaRange(chroma_range_table[color].floor, chroma_range_table[color].ceil)
	vsw.SetSubMode(SUB_MODE_CHROMAKEY)
}
