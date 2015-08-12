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
type Chroma_range struct {
	floor, ceil int
}

var chroma_range_table = []Chroma_range{
	chroma_range_entry(0, 0, 8-1, 8-2),
	chroma_range_entry(8+1, 0, 15, 8-2),
	chroma_range_entry(8+1, 8+2, 15, 15),
	chroma_range_entry(0, 8+2, 8-1, 15),
}

func chroma_range_entry(u0, v0, u1, v1 int) Chroma_range {
	return Chroma_range{
		(u0)<<4 | 1<<12 | (v0)<<20,
		(u1)<<4 | 14<<12 | (v1)<<20,
	}
}

// SetChromaRange sets color range of chroma key
func (vsw Vsw) SetChromaRange(r Chroma_range) {
	cmd := []uint32{SW_ID_SetChromaRange, uint32(r.floor), uint32(r.ceil)}
	send(vsw.conn, cmd)
}

// SetChromaKey sets chroma key in predefined color
func (vsw Vsw) SetChromaKey(color int) {
	if 0 > color || color >= CHROMA_KEY_NUM {
		return
	}
	vsw.SetChromaRange(chroma_range_table[color])
	vsw.SetSubMode(SUB_MODE_CHROMAKEY)
}
