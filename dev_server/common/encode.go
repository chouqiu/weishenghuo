package common

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"github.com/golang/glog"
)

const (
	coding_key = "artzwei12798akljzmknm.ahkjkljl;k"
	iv_key     = "artzwei12798aklj"
)

func UrlEncode(src []byte) ([]byte, error) {
	return aesEncrypt(src)
}

func UrlDecode(src []byte) ([]byte, error) {
	return aesDecrypt(src)
}

func Base64Encode(src []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(src))
}

func Base64Decode(src []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(src))
}

func aesEncrypt(plaintext []byte) ([]byte, error) {
	// 创建加密算法aes
	c, err := aes.NewCipher([]byte(coding_key))
	if err != nil {
		glog.Errorf("Error: NewCipher(%d bytes) = %s", len(coding_key), err)
		return nil, err
	}

	cfb := cipher.NewCFBEncrypter(c, []byte(iv_key))
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)

	return ciphertext, nil
}

func aesDecrypt(ciphertext []byte) ([]byte, error) {
	// 创建加密算法aes
	c, err := aes.NewCipher([]byte(coding_key))
	if err != nil {
		glog.Errorf("Error: NewCipher(%d bytes) = %s", len(coding_key), err)
		return nil, err
	}

	cfbdec := cipher.NewCFBDecrypter(c, []byte(iv_key))
	plaintextCopy := make([]byte, len(ciphertext))
	cfbdec.XORKeyStream(plaintextCopy, ciphertext)

	return plaintextCopy, nil
}
