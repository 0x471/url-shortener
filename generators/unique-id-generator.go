package generators

import (
	"strconv"
	"hash/adler32"
)

func GenerateUniqueID(url string) string {
	return strconv.FormatUint(uint64(adler32.Checksum([]byte(url))), 16)
}
