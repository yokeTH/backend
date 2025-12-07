package idgenerator

import "strconv"

func strconvU16(v uint16) string {
	return strconv.FormatUint(uint64(v), 10)
}
