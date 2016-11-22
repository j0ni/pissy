package db

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"errors"
	"fmt"
)

const SaltSize = 16
const IVSize = aes.BlockSize
const AESKeySize = 32

func checkHMAC(value, expected, key []byte) bool {
	mac := hmac.New(sha512.New, key)
	mac.Write(value)
	return hmac.Equal(expected, mac.Sum(nil))
}

func generateHMAC(value, key []byte) []byte {
	mac := hmac.New(sha512.New, key)
	mac.Write(value)
	return mac.Sum(nil)
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

func secureRandomBytes(size int) []byte {
	bs := make([]byte, size)
	_, err := rand.Read(bs)
	if err != nil {
		panic(err)
	}
	return bs
}

func decrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New(fmt.Sprintf("ciphertext is too short (%d bytes, %d required)",
			len(ciphertext), aes.BlockSize))
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)
	return stripPadding(plaintext)
}

func encrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plaintext = pad(plaintext)
	ciphertext := make([]byte, len(plaintext))
	iv := secureRandomBytes(IVSize)
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)
	return append(iv, ciphertext...), nil
}
