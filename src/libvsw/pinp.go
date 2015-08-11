package libvsw

// SetSubMode sets the sub screen mode
func (vsw Vsw) SetSubMode(mode int) {
	cmd := []uint32{SW_ID_SetSubMode, uint32(mode)}
	send(vsw.conn, cmd)
}

// PinpGeometry sets the scale and crop of the sub screen.
func (vsw Vsw) PinpGeometry(scale_x, scale_y, scale_w, scale_h, crop_x, crop_y, crop_w, crop_h int) {
	cmd := []uint32{SW_ID_SetPinpGeometry,
		uint32(scale_x), uint32(scale_y), uint32(scale_w), uint32(scale_h),
		uint32(crop_x), uint32(crop_y), uint32(crop_w), uint32(crop_h)}
	send(vsw.conn, cmd)
}

// PinBorder sets color and width of border of the sub screen.
func (vsw Vsw) PinpBorder(color int, width int) {
	cmd := []uint32{SW_ID_SetPinpBorder, uint32(color), uint32(width)}
	send(vsw.conn, cmd)
}
