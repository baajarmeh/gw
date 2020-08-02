package gw

import (
	"crypto/cipher"
	"github.com/oceanho/gw/utils/secure"
	"sync"
)

type ICrypto interface {
	Hash() ICryptoHash
	Protect() ICryptoProtect
}

type ICryptoHash interface {
	Hash(dst, src []byte) error
}

type ICryptoProtect interface {
	Encrypt(dst, src []byte) error
	Decrypt(dst, src []byte) error
}

var (
	cryptoOnce        sync.Once
	defaultCryptoImpl *DefaultCryptoImpl
)

type DefaultCryptoImpl struct {
	secret string
	salt   string
}

func DefaultCrypto(secret, salt string) ICrypto {
	cryptoOnce.Do(func() {
		defaultCryptoImpl = &DefaultCryptoImpl{
			secret: secret,
			salt:   salt,
		}
	})
	return defaultCryptoImpl
}

func (d DefaultCryptoImpl) Hash() ICryptoHash {
	return DefaultCryptoHashSha256(d.salt)
}

func (d DefaultCryptoImpl) Protect() ICryptoProtect {
	return DefaultCryptoProtectAES(d.secret)
}

type DefaultCryptoProtectAESImpl struct {
	key       string
	encrypter cipher.Stream
	decrypter cipher.Stream
}

var (
	onceAES                     sync.Once
	onceSha256                  sync.Once
	defaultCryptoProtectAES     *DefaultCryptoProtectAESImpl
	defaultCryptoHashSha256Impl *DefaultCryptoHashSha256Impl
)

func DefaultCryptoProtectAES(key string) *DefaultCryptoProtectAESImpl {
	onceAES.Do(func() {
		b := secure.AesBlock(key)
		defaultCryptoProtectAES = &DefaultCryptoProtectAESImpl{
			key:       key,
			encrypter: secure.AesEncryptCFB(key, b),
			decrypter: secure.AesDecryptCFB(key, b),
		}
	})
	return defaultCryptoProtectAES
}

func (d *DefaultCryptoProtectAESImpl) Encrypt(dst, src []byte) error {
	// d.encrypter.XORKeyStream(dst, src)
	b := secure.AesBlock(d.key)
	secure.AesEncryptCFB(d.key, b).XORKeyStream(dst, src)
	return nil
}

func (d DefaultCryptoProtectAESImpl) Decrypt(dst, src []byte) error {
	//
	// FIXME(Ocean): Block must A new for every call. why ?
	// d.decrypter.XORKeyStream(dst, src)
	//
	b := secure.AesBlock(d.key)
	secure.AesDecryptCFB(d.key, b).XORKeyStream(dst, src)
	return nil
}

type DefaultCryptoHashSha256Impl struct {
	salt string
}

func DefaultCryptoHashSha256(salt string) *DefaultCryptoHashSha256Impl {
	onceSha256.Do(func() {
		defaultCryptoHashSha256Impl = &DefaultCryptoHashSha256Impl{
			salt: salt,
		}
	})
	return defaultCryptoHashSha256Impl
}

func (d *DefaultCryptoHashSha256Impl) Hash(dst, src []byte) error {
	dst = []byte(secure.Sha256(src))
	return nil
}
