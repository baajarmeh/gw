package secure

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAES_Key16_Test100Counter(t *testing.T) {
	key := "AWpqKjLdLheNnPo3"
	block := AesBlock(key)
	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("dE7axViJTcYgygphmWMNsVTciHXar7y,%d", i)
		encrypted := make([]byte, len(msg))
		decrypted := make([]byte, len(msg))
		a := AesEncryptCFB(key, block)
		EncryptAES(encrypted, []byte(msg), a)
		b := AesDecryptCFB(key, block)
		DecryptAES(decrypted, encrypted, b)
		assert.Equal(t, string(decrypted), msg)
	}
}

func BenchmarkAES_Key16(b *testing.B) {
	key := "AWpqKjLdLheNnPo3"
	block := AesBlock(key)
	for n := 0; n < b.N; n++ {
		msg := fmt.Sprintf("dE7axViJTcYgygphmWMNsVTciHXar7y,%d", n)
		encrypted := make([]byte, len(msg))
		decrypted := make([]byte, len(msg))
		EncryptAES(encrypted, []byte(msg), AesEncryptCFB(key, block))
		DecryptAES(decrypted, encrypted, AesDecryptCFB(key, block))
	}
}

func BenchmarkAES_Key24(b *testing.B) {
	key := "pjLa3yYmeKwjLLnYYnpazvLL"
	block := AesBlock(key)
	for n := 0; n < b.N; n++ {
		msg := fmt.Sprintf("dE7axViJTcYgygphmWMNsVTciHXar7y,%d", n)
		encrypted := make([]byte, len(msg))
		decrypted := make([]byte, len(msg))
		EncryptAES(encrypted, []byte(msg), AesEncryptCFB(key, block))
		DecryptAES(decrypted, encrypted, AesDecryptCFB(key, block))
	}
}

func BenchmarkAES_Key32(b *testing.B) {
	key := "bxrXXxxdvVfgbhshhxptavcrKd3JwUjk"
	block := AesBlock(key)
	for n := 0; n < b.N; n++ {
		msg := fmt.Sprintf("dE7axViJTcYgygphmWMNsVTciHXar7y,%d", n)
		encrypted := make([]byte, len(msg))
		decrypted := make([]byte, len(msg))
		EncryptAES(encrypted, []byte(msg), AesEncryptCFB(key, block))
		DecryptAES(decrypted, encrypted, AesDecryptCFB(key, block))
	}
}
