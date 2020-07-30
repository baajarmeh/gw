package secure

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

//
// ref:
// https://blog.csdn.net/xiaohu50/article/details/51682849
//
func EncryptAES(dst, src []byte, aesEncrypter cipher.Stream) {
	aesEncrypter.XORKeyStream(dst, src)
}

func DecryptAES(dst, src []byte, aesDecrypter cipher.Stream) {
	aesDecrypter.XORKeyStream(dst, src)
}

func AesBlock(key string) cipher.Block {
	b, e := aes.NewCipher([]byte(key))
	if e != nil {
		panic(fmt.Sprintf("aes.NewCipher(...), err:%v", e))
	}
	return b
}

func AesEncryptCFB(key string, block cipher.Block) cipher.Stream {
	var iv = []byte(key)[:aes.BlockSize]
	return cipher.NewCFBEncrypter(block, iv)
}

func AesDecryptCFB(key string, block cipher.Block) cipher.Stream {
	kl := len(key)
	if kl != 16 && kl != 24 && kl != 32 {
		panic("key size only supports (16,24,32).")
	}
	var iv = []byte(key)[:aes.BlockSize]
	return cipher.NewCFBDecrypter(block, iv)
}
