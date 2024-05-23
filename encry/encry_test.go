package encry

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"testing"
)

/**
@file:
@author: levi.Tang
@time: 2024/4/15 10:35
@description:
**/

func TestDecrypt(t *testing.T) {
	src := "FA52OfBQVlSP1w=="
	key := []byte("62me6ia5bf46wo84lk7qag3dmhmdpul3")

	decode, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		t.Errorf("init db connect failed, decode, err: %s", err.Error())
		os.Exit(1)
	}

	var iv = []byte(key)[:aes.BlockSize]
	decrypted := make([]byte, len(decode))
	var aesBlockDecrypter cipher.Block
	aesBlockDecrypter, err = aes.NewCipher([]byte(key))
	if err != nil {
		t.Errorf("init db connect failed, decrypt, err: %s ", err.Error())
		os.Exit(1)
	}
	aesDecrypter := cipher.NewCFBDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.XORKeyStream(decrypted, decode)
	fmt.Println(string(decrypted))
}

func TestEncrypt(t *testing.T) {
	src := []byte("myviu_4359")
	key := []byte("62me6ia5bf46wo84lk7qag3dmhmdpul3")

	var aesBlockDecrypter cipher.Block
	aesBlockDecrypter, err := aes.NewCipher([]byte(key))
	if err != nil {
		t.Errorf("init db connect failed, decrypt, err: %s ", err.Error())
		os.Exit(1)
	}
	ciphertext := make([]byte, aes.BlockSize+len(src))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	mode := cipher.NewCFBEncrypter(aesBlockDecrypter, iv)
	mode.XORKeyStream(ciphertext, src)
	// 记住密文必须经过认证是很重要的
	// （即通过使用crypto/hmac）以及为了加密而被加密
	// 保持secure。

	fmt.Printf("%x\n", ciphertext)
	var encry []byte
	aesBlockDecrypter.Encrypt(encry, []byte(src))
	fmt.Println(string(encry))
}
