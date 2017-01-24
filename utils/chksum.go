package utils

/* WARNING: len1 MUST be an EVEN number */
func DualbufferChecksum(inbuf1 []byte, inbuf2 []byte) uint16 {

	sum := checksum_adder(0, inbuf1)
	sum = checksum_adder(sum, inbuf2)
	return checksum_finalize(sum)
}

/**
 * Calculate checksum of a given string
 */
func Checksum_(inbuf []byte) uint16 {
	sum := checksum_adder(0, inbuf)
	return checksum_finalize(sum)
}
func Checksum(data []byte) uint16 {
	var chksum uint32

	var lsb uint16
	var msb uint16

	// 32-bit sum (2's complement sum of 16 bits with carry)
	for i := 0; i < len(data)-1; i += 2 {
		msb = uint16(data[i])
		lsb = uint16(data[i+1])
		chksum += uint32(lsb + (msb << 8))
	}

	// 1's complement 16-bit sum via "end arround carry" of 2's complement
	chksum = ((chksum >> 16) & 0xFFFF) + (chksum & 0xFFFF)

	return uint16(0xFFFF & (^chksum))
}

func VerifyChecksum(data []byte) bool {
	return Checksum(data) == 0
}

func checksum_adder(sum uint32, data []byte) uint32 {

	var lsb uint16
	var msb uint16
	lenth := len(data)
	if lenth&0x01 > 0 {
		lenth--
		//    sum += (((uint8_t *)data)[len]) << 8;
		sum += uint32(data[lenth])
	}

	// 32-bit sum (2's complement sum of 16 bits with carry)
	for i := 0; i < len(data)-1; i += 2 {
		msb = uint16(data[i])
		lsb = uint16(data[i+1])
		sum += uint32(lsb + (msb << 8))
	}
	return sum
}

func checksum_finalize(sum uint32) uint16 {
	for (sum >> 16) > 0 { /* a second carry is possible! */
		sum = (sum & 0x0000FFFF) + (sum >> 16)
	}
	return ShortBe(uint16(^sum))
}
