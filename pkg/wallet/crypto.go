package wallet

import (
	"unsafe"
	"github.com/awnumar/memguard"
	"crypto/aes"
	"crypto/sha512"
	"crypto/cipher"
)
// Encrypts plaintext using the masterkey and a password
func (m *MKey) Encrypt(p string, b ...[]byte) (r [][]byte, err error) {
	return m.encDec(true, p, b...)
}

// Decrypts a ciphertext using the masterkey and password
func (m *MKey) Decrypt(p string, b ...[]byte) (r [][]byte, err error) {
	return m.encDec(false, p, b...)
}
func (m *MKey) encDec(enc bool, p string, b ...[]byte) (r [][]byte, err error) {
	source, iv, err := m.DeriveCipher(p)
	if block, err := aes.NewCipher(source.Buffer()); err != nil {
		return nil, err
	} else {
		for i := range b {
			if enc {
				mode := cipher.NewCBCEncrypter(block, iv)
				mode.CryptBlocks(r[i], b[i])
			} else {
				mode := cipher.NewCBCDecrypter(block, iv)
				mode.CryptBlocks(r[i], b[i])
			}
		}
		for i := range source.Buffer() { source.Buffer()[i] = 0 }
		for i := range b { 
			for j := range b[i] {
				b[i][j] = 0
			}
		}
	}	
	return
}

func (m *MKey) DeriveCipher(pass string) (k *memguard.LockedBuffer, iv []byte, err error) {
	var seed *memguard.LockedBuffer
	pLen, sLen := len(pass), len(m.Salt)
	if pLen + sLen > 64 {
		sLen = 128 - pLen - sLen
		seed, err = memguard.NewMutable(128)
	} else {
		seed, err = memguard.NewMutable(64)
	}
	if err != nil {
		return
	}
	buf := seed.Buffer()
	for i := range pass {
		buf[i] = pass[i]
	}
	for i := range m.Salt {
		buf[i+pLen] = m.Salt[i]
	}
	// PKCS#5 padding - pad byte is number of pad bytes max 64
	pad := len(buf) - pLen - sLen
	for i := pLen + sLen; i < len(buf); i++ {
		buf[i] = byte(pad)
	}
	var source *[64]byte
	l, err := memguard.NewMutable(64)
	source = (*[64]byte)(unsafe.Pointer(&l.Buffer()[0]))
	*source = sha512.Sum512(buf)
	for i := 0; i < int(m.Iterations-1); i++ {
		*source = sha512.Sum512(l.Buffer())
	}
	k, err = memguard.NewMutable(64)
	for i := range k.Buffer() {
		k.Buffer()[i] = source[i]
	}
	ckey, ivb, err := memguard.Split(k, 32)
	block, err := aes.NewCipher(ckey.Buffer())
	iv = ivb.Buffer()[:block.BlockSize()]
   mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(k.Buffer(), m.EncryptedKey)
	l.Destroy()
	seed.Destroy()
	return
}

