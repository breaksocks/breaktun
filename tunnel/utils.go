package tunnel

func WriteN2(bs []byte, offset int, n uint16) {
	bs[offset] = byte(n >> 8)
	bs[offset+1] = byte(n & 0xFF)
}

func ReadN2(bs []byte, offset int) uint16 {
	return (uint16(bs[offset]) << 8) | uint16(bs[offset+1])
}

func WriteN4(bs []byte, offset int, n uint32) {
	bs[offset] = byte(n >> 24)
	bs[offset+1] = byte(n >> 16)
	bs[offset+2] = byte(n >> 8)
	bs[offset+3] = byte(n)
}

func ReadN4(bs []byte, offset int) uint32 {
	var n uint32
	n |= uint32(bs[offset]) << 24
	n |= uint32(bs[offset+1]) << 16
	n |= uint32(bs[offset+2]) << 8
	n |= uint32(bs[offset+3])
	return n
}
