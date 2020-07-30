package secure

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
)

//
// ref:
// https://blog.csdn.net/wade3015/article/details/84454836
//
func RsaEncrypt(pub *rsa.PublicKey, origData []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

func RsaDecrypt(pri *rsa.PrivateKey, cipherText []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, pri, cipherText)
}

func LoadPrivateKeyFile(filepath string) *rsa.PrivateKey {
	k, e := ioutil.ReadFile(filepath)
	if e != nil {
		panic(e)
	}
	block, _ := pem.Decode(k)
	if block == nil {
		panic(e)
	}
	pri, e := x509.ParsePKCS1PrivateKey(block.Bytes)
	if e != nil {
		panic(e)
	}
	return pri
}

func LoadPublicKeyFile(filepath string) *rsa.PublicKey {
	p, e := ioutil.ReadFile(filepath)
	if e != nil {
		panic(e)
	}
	block, _ := pem.Decode(p)
	if block == nil {
		panic(e)
	}
	pub, e := x509.ParsePKIXPublicKey(block.Bytes)
	if e != nil {
		panic(e)
	}
	return pub.(*rsa.PublicKey)
}
