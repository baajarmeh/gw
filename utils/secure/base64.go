package secure

import (
	"encoding/base64"
	"strings"
)

//
// refer:
// https://github.com/dvsekhvalnov/jose2go/blob/master/base64url/base64url.go
//

// DecodeBase64URL decodes base64url string to byte array
func DecodeBase64URL(data string) ([]byte, error) {
	data = strings.Replace(data, "-", "+", -1) // 62nd char of encoding
	data = strings.Replace(data, "_", "/", -1) // 63rd char of encoding

	switch len(data) % 4 { // Pad with trailing '='s
	case 0: // no padding
	case 2:
		data += "==" // 2 pad chars
	case 3:
		data += "=" // 1 pad char
	}
	return base64.StdEncoding.DecodeString(data)
}

// EncodeBase64URL encodes given byte array to base64url string
func EncodeBase64URL(data []byte) string {
	result := base64.StdEncoding.EncodeToString(data)
	result = strings.Replace(result, "+", "-", -1) // 62nd char of encoding
	result = strings.Replace(result, "/", "_", -1) // 63rd char of encoding
	result = strings.Replace(result, "=", "", -1)  // Remove any trailing '='s
	return result
}
