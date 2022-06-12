package cryptx

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/bcrypt"

	"golang.org/x/crypto/scrypt"
)

// EncryptSalt 加盐加密，相同输入会有相同的结果
func EncryptSalt(salt, src string) string {
	dk, _ := scrypt.Key([]byte(src), []byte(salt), 32768, 8, 1, 32)
	return fmt.Sprintf("%x", dk)
}

func CompareSalt(salt, encrypt, input string) bool {
	return EncryptSalt(salt, input) == encrypt
}

// EncryptRandom 相同输入每次会生成不一样的结果
func EncryptRandom(src string) string {
	dk, _ := bcrypt.GenerateFromPassword([]byte(src), bcrypt.DefaultCost)
	return string(dk)
}

func CompareRandom(encrypt, input string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encrypt), []byte(input)) == nil
}

// EncryptSHA256 使用 SHA256 加密密码存储到数据库中
func EncryptSHA256(src string) string {
	hash := sha256.New()
	hash.Write([]byte(src))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func CompareSHA256(encrypt, src string) bool {
	return EncryptSHA256(src) == encrypt
}

// EncryptAES AES 加密
func EncryptAES(key, origData []byte) []byte {
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()
	origData = pkcS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	encrypted := make([]byte, len(origData))
	blockMode.CryptBlocks(encrypted, origData)
	return encrypted
}

// DecryptAES AES 解密
func DecryptAES(key, encrypt []byte) []byte {
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(encrypt))
	blockMode.CryptBlocks(origData, encrypt)
	origData = pkcS7UnPadding(origData)
	return origData
}

func pkcS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// EncryptXOR 异或加密
func EncryptXOR(key, origData []byte) []byte {
	origData = pkcS7Padding(origData, len(key))
	for i := 0; i < len(origData); i++ {
		origData[i] ^= key[i%len(key)]
	}
	return origData
}

// DecryptXOR 异或解密
func DecryptXOR(key, encrypt []byte) []byte {
	for i := 0; i < len(encrypt); i++ {
		encrypt[i] ^= key[i%len(key)]
	}
	return pkcS7UnPadding(encrypt)
}
