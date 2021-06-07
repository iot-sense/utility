package utility

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

var key = []byte(G_CONFIGER.GetString("aesKey"))

//AesEncryptCFB func
func AesEncryptCFB(input string) (output string) {
	origData := []byte(input)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	encrypted := make([]byte, aes.BlockSize+len(origData))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], origData)
	output = base64.StdEncoding.EncodeToString(encrypted)
	return output
}

//AesDecryptCFB func
func AesDecryptCFB(input string) (output string) {
	encrypted, _ := base64.StdEncoding.DecodeString(input)
	block, _ := aes.NewCipher(key)
	if len(encrypted) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]
	decrypted := make([]byte, len(encrypted))
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(decrypted, encrypted)
	output = string(decrypted)
	return output
}

//TestAes func
func TestAes() {
	data := "password"
	Logger.Info("Origin：", data)

	encrypted := AesEncryptCFB(data)
	Logger.Info("Encrypted：", encrypted)

	decrypted := AesDecryptCFB(encrypted)
	Logger.Info("Decrypted：", decrypted)
}
