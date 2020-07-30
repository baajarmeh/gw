package secure

import (
	"crypto/rsa"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	pri *rsa.PrivateKey
	pub *rsa.PublicKey
)

func init() {
	pri = LoadPrivateKeyFile("files/gw.key")
	pub = LoadPublicKeyFile("files/gw.pem")
}

func TestRsa(t *testing.T) {
	str := "I am oceanho."
	byt, err := RsaEncrypt(pub, []byte(str))
	t.Logf("encrypt,err:%v, data is: %v", err, byt)

	ori, err := RsaDecrypt(pri, byt)
	t.Logf("decrypt,err:%v, data is: %v", err, ori)

	assert.Equal(t, string(ori), str)
}

func BenchmarkRsa(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		str := "I am oceanho."
		byt, _ := RsaEncrypt(pub, []byte(str))
		ori, _ := RsaDecrypt(pri, byt)
		assert.Equal(b, string(ori), str)
	}
}
