package main

import (
	"encoding/base64"
	"fmt"

	"github.com/ph9/go-lib/crypto"
)

func encryptByConfig(plainText string) (string, error) {
	if plainText == "" {
		return plainText, nil
	}

	config := configInstance()
	privateKey := []byte(config.PrivateKey)
	bIv := []byte(config.IvKey)
	b, err := crypto.EncryptWithTripleDes(privateKey, bIv, []byte(plainText))

	if err != nil {
		fmt.Println(err)
		return plainText, err
	}

	result := base64.URLEncoding.EncodeToString(b)
	return result, nil
}

func decrypt(base64Text, pKey, ivKey string) (string, error) {
	bPKey := []byte(pKey)
	bIv := []byte(ivKey)
	b, err := crypto.DecryptWithTripleDes(bPKey, bIv, []byte(base64Text))
	if err != nil {
		return base64Text, err
	}

	return string(b), err
}
