package secretutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Sha1(content string) string {
	h := sha1.New()
	_, _ = h.Write([]byte(content))
	binaryData := h.Sum(nil)
	return hex.EncodeToString(binaryData)

}

// Encrypt 加密 (支持自定义 IV)
func AesEncrypt(text, key, iv []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(iv) != aes.BlockSize {
		return "", errors.New("IV length must be equal to AES block size")
	}

	plainText := padding(text, aes.BlockSize)
	cipherText := make([]byte, len(plainText))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, plainText)

	// 将加密结果转换为 Base64 编码
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt 解密 (支持自定义 IV)
func AesDecrypt(cipherTextBase64, key, iv []byte) (string, error) {
	if len(cipherTextBase64) <= 0 {
		return "", nil
	}

	// 解码 Base64 加密内容
	cipherText, err := base64.StdEncoding.DecodeString(string(cipherTextBase64))
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(iv) != aes.BlockSize {
		return "", errors.New("IV length must be equal to AES block size")
	}

	if len(cipherText)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	plainText := unPadding(cipherText)
	return string(plainText), nil
}

// Padding 用于填充明文，保证它的长度是块大小的整数倍
func padding(plainText []byte, blockSize int) []byte {
	padding := blockSize - len(plainText)%blockSize
	padText := make([]byte, padding)
	for i := range padText {
		padText[i] = byte(padding)
	}
	return append(plainText, padText...)
}

// UnPadding 去除填充数据
func unPadding(plainText []byte) []byte {
	length := len(plainText)
	unpadding := int(plainText[length-1])
	return plainText[:(length - unpadding)]
}
