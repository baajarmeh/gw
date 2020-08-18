package gw

import (
	"crypto/cipher"
	"fmt"
	"github.com/oceanho/gw/utils/secure"
	"sync"
)

type ICrypto interface {
	Hash() ICryptoHash
	Protect() ICryptoProtect
	Password() IPasswordSigner
}

type ICryptoHash interface {
	Hash(dst, src []byte) error
}

type ICryptoProtect interface {
	Encrypt(dst, src []byte) error
	Decrypt(dst, src []byte) error
}

type IPasswordSigner interface {
	Sign(plainPassword string) string
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

func (d DefaultCryptoImpl) Password() IPasswordSigner {
	return DefaultPasswordSignerMd5(d.salt)
}

type DefaultCryptoProtectAESImpl struct {
	key       string
	encrypter cipher.Stream
	decrypter cipher.Stream
}

var (
	onceAES                       sync.Once
	onceSha256                    sync.Once
	onceMd5                       sync.Once
	defaultCryptoProtectAES       *DefaultCryptoProtectAESImpl
	defaultCryptoHashSha256Impl   *DefaultCryptoHashSha256Impl
	defaultPasswordProtectMd5Impl *DefaultPasswordSignerMd5Impl
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

func DefaultPasswordSignerMd5(salt string) *DefaultPasswordSignerMd5Impl {
	onceMd5.Do(func() {
		defaultPasswordProtectMd5Impl = &DefaultPasswordSignerMd5Impl{
			salt: salt,
		}
	})
	return defaultPasswordProtectMd5Impl
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

type DefaultPasswordSignerMd5Impl struct {
	salt string
}

func (d DefaultPasswordSignerMd5Impl) str(ori string) string {
	return fmt.Sprintf("%s,./;@#$%s,(*)(^$#%s,./;@#$%s,(", d.salt, ori, ori, d.salt)
}

func (d DefaultPasswordSignerMd5Impl) Sign(plainPassword string) string {
	return secure.Md5Str(d.str(plainPassword))
}
