package db

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/gob"
	"errors"
	"fmt"

	"golang.org/x/crypto/pbkdf2"

	uuid "github.com/satori/go.uuid"
)

type EncryptionKey struct {
	Uuid             uuid.UUID
	EncryptedKey     []byte
	DecryptedKey     []byte
	EncryptedKeyHMAC []byte
	Salt             []byte
	Iterations       int
}

const SaltSize = 16
const IVSize = aes.BlockSize
const AESKeySize = 32

func (key *EncryptionKey) Load(dir, fileName string) error {
	buf, err := loadFile(dir, fileName)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(buf)
	return dec.Decode(key)
}

func (key EncryptionKey) Save(dir, fileName string) error {
	var buf bytes.Buffer
	// its a copy
	key.DecryptedKey = nil
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(key); err != nil {
		return err
	}
	return saveFile(dir, fileName, buf)
}

func pad(text []byte) []byte {
	buf := bytes.NewBuffer(text)
	rem := len(text) % aes.BlockSize
	if rem == 0 {
		rem = aes.BlockSize
	}
	for i := 0; i < rem; i++ {
		err := buf.WriteByte(byte(rem))
		if err != nil {
			panic(err)
		}
	}
	return buf.Bytes()
}

func stripPadding(text []byte) ([]byte, error) {
	lastByte := text[len(text)-1]
	textLen := len(text) - int(lastByte)
	padding := text[textLen:]
	// don't do this because
	// https://www.nccgroup.trust/us/about-us/newsroom-and-events/blog/2009/july/if-youre-typing-the-letters-a-e-s-into-your-code-youre-doing-it-wrong/
	// so is the right thing to do to just silently continue and
	// generate garbage? I guess so. Useful for debugging so TODO FIXME HOLYSHIT
	for _, v := range padding {
		if v != lastByte {
			return nil, errors.New("Padding integrity failure")
		}
	}
	return text[:textLen], nil
}

func (key EncryptionKey) deriveKey(passphrase string, salt []byte) []byte {
	return pbkdf2.Key([]byte(passphrase), salt, key.Iterations, 32, sha512.New)
}

func secureRandomBytes(size int) []byte {
	bs := make([]byte, size)
	_, err := rand.Read(bs)
	if err != nil {
		panic(err)
	}
	return bs
}

func (key *EncryptionKey) updateEncryptedKeyHMAC(macKey []byte) {
	mac := hmac.New(sha512.New, macKey)
	mac.Write(key.EncryptedKey)
	key.EncryptedKeyHMAC = mac.Sum(nil)
}

func (key EncryptionKey) checkValidationHMAC(macKey []byte) bool {
	mac := hmac.New(sha512.New, macKey)
	mac.Write(key.EncryptedKey)
	return hmac.Equal(key.EncryptedKeyHMAC, mac.Sum(nil))
}

func (key *EncryptionKey) DecryptKey(passphrase string) error {
	dk := key.deriveKey(passphrase, key.Salt)
	// not sure if it's ok to reuse this key
	if !key.checkValidationHMAC(dk) {
		return errors.New("HMAC validation failed")
	}
	block, err := aes.NewCipher(dk)
	if err != nil {
		return err
	}
	if len(key.EncryptedKey) < aes.BlockSize {
		return errors.New(fmt.Sprintf("EncryptedKey is too short (%d bytes, %d required)",
			len(key.EncryptedKey), aes.BlockSize))
	}
	iv := key.EncryptedKey[:aes.BlockSize]
	ciphertext := key.EncryptedKey[aes.BlockSize:]
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)
	key.DecryptedKey, err = stripPadding(plaintext)
	return err
}

func (key *EncryptionKey) EncryptKey(passphrase string) error {
	key.Salt = secureRandomBytes(SaltSize)
	dk := key.deriveKey(passphrase, key.Salt)
	block, err := aes.NewCipher(dk)
	if err != nil {
		return err
	}
	padded := pad(key.DecryptedKey)
	ciphertext := make([]byte, len(padded))
	iv := secureRandomBytes(IVSize)
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padded)
	key.EncryptedKey = append(iv, ciphertext...)
	// not sure if it's ok to reuse this key
	key.updateEncryptedKeyHMAC(dk)
	return nil
}

func (key *EncryptionKey) GenerateKey() {
	key.DecryptedKey = secureRandomBytes(AESKeySize)
}
