package cryptox

import (
	"crypto/aes"
	"crypto/cipher"
)

func AESEncodeByte(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return data, err
	}
	var iv = key[:aes.BlockSize]
	blockMode := cipher.NewCFBEncrypter(block, iv)
	dest := make([]byte, len(string(data)))
	blockMode.XORKeyStream(dest, data)
	return dest, nil
}
func AESDecodeByte(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return data, err
	}
	var iv = key[:aes.BlockSize]
	blockMode := cipher.NewCFBDecrypter(block, iv)
	dest := make([]byte, len(string(data)))
	blockMode.XORKeyStream(dest, data)
	return dest, nil
}
func AESEncode(data string, key string) (string, error) {
	out, err := AESEncodeByte([]byte(data), []byte(key))
	return string(Base64EncodeByte(out)), err
}
func AESDecode(data string, key string) (string, error) {
	d, e := Base64DecodeByte([]byte(data))
	if e != nil {
		return data, e
	}
	out, err := AESDecodeByte(d, []byte(key))
	return string(out), err
}

// AES加密
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// AES解密
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}
