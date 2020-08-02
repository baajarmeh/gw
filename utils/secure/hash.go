package secure

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
)

//
// ref
// https://stackoverflow.com/questions/2377881/how-to-get-a-md5-hash-from-a-string-in-golang
//
// Md5 ...
func Md5(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}

func Md5Str(data string) string {
	return Md5([]byte(data))
}

func Sha256(data []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(data))
}

func Sha256Str(data string) string {
	return Sha256([]byte(data))
}
