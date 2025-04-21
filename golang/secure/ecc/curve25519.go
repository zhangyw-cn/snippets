package ecc

import (
	"crypto/rand"
	"io"

	"golang.org/x/crypto/curve25519"
)

// 生成 X25519 密钥对
func generateKeyPair() (privateKey, publicKey [32]byte, err error) {
	if _, err = io.ReadFull(rand.Reader, privateKey[:]); err != nil {
		return
	}
	curve25519.ScalarBaseMult(&publicKey, &privateKey)
	return
}

// 使用自己的私钥 + 其他人的公钥 生成共享密钥
func computeSharedSecret(privateKey, publicKey [32]byte) (sharedKey [32]byte) {
	curve25519.ScalarMult(&sharedKey, &privateKey, &publicKey)
	return
}
