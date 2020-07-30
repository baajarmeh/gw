package secure

import (
	"crypto/cipher"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	aesEncrypter cipher.Stream
	aesDecrypter cipher.Stream
)

func init() {
	var key = "km399jkcqYbqLTKL"
	block := AesBlock(key)
	aesEncrypter = AesEncryptCFB(key, block)
	aesDecrypter = AesDecryptCFB(key, block)
}

func TestAESWithSameCipher(t *testing.T) {
	const msg = "dE7axViJTcYgygphmWMNsVTciHXar7yy"

	encrypted := make([]byte, len(msg))
	EncryptAES(encrypted, []byte(msg), aesEncrypter)

	decrypted := make([]byte, len(msg))
	DecryptAES(decrypted, []byte(encrypted), aesDecrypter)
	assert.Equal(t, string(decrypted), msg)
}

func BenchmarkAESWithSameCipher(b *testing.B) {
	const msg = "dE7axViJTcYgygphmWMNsVTciHXar7yy"
	for n := 0; n < b.N; n++ {
		encrypted := make([]byte, len(msg))
		EncryptAES(encrypted, []byte(msg), aesEncrypter)
		decrypted := make([]byte, len(msg))
		DecryptAES(decrypted, []byte(encrypted), aesDecrypter)
	}
}
