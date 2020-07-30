package gw

import (
	"crypto/cipher"
	"github.com/oceanho/gw/utils/secure"
	"sync"
)

type ICrypto interface {
	Encrypt(dst, src []byte) error
	Decrypt(dst, src []byte) error
}

type DefaultCryptoAesImpl struct {
	key       string
	encrypter cipher.Stream
	decrypter cipher.Stream
}

var (
	once             sync.Once
	defaultCryptoAes *DefaultCryptoAesImpl
)

func DefaultCryptoAes(key string) *DefaultCryptoAesImpl {
	once.Do(func() {
		b := secure.AesBlock(key)
		defaultCryptoAes = &DefaultCryptoAesImpl{
			key:       key,
			encrypter: secure.AesEncryptCFB(key, b),
			decrypter: secure.AesDecryptCFB(key, b),
		}
	})
	return defaultCryptoAes
}

func (d *DefaultCryptoAesImpl) Encrypt(dst, src []byte) error {
	// d.encrypter.XORKeyStream(dst, src)
	b := secure.AesBlock(d.key)
	secure.AesEncryptCFB(d.key, b).XORKeyStream(dst, src)
	return nil
}

func (d DefaultCryptoAesImpl) Decrypt(dst, src []byte) error {
	//
	// FIXME(Ocean): Block must A new for every call. why ?
	// d.decrypter.XORKeyStream(dst, src)
	//
	b := secure.AesBlock(d.key)
	secure.AesDecryptCFB(d.key, b).XORKeyStream(dst, src)
	return nil
}
