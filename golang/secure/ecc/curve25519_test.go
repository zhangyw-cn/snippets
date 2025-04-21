package ecc

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"
)

// 密钥交换（X25519）
func TestCurve25519_X25519(t *testing.T) {
	// Alice 和 Bob 生成密钥对
	alicePriv, alicePub, _ := generateKeyPair()
	bobPriv, bobPub, _ := generateKeyPair()

	// 计算共享密钥
	sharedAlice := computeSharedSecret(alicePriv, bobPub)
	sharedBob := computeSharedSecret(bobPriv, alicePub)

	// 输出结果（应相同）
	println("Shared Key Match:", sharedAlice == sharedBob)
}

// 数字签名（Ed25519）
func TestCurve25519_Ed25519(t *testing.T) {
	// 生成密钥对
	pubKey, privKey, _ := ed25519.GenerateKey(rand.Reader)

	// 签名消息
	message := []byte("Hello, Ed25519!")
	signature := ed25519.Sign(privKey, message)

	// 验证签名
	isValid := ed25519.Verify(pubKey, message, signature)
	println("Signature Valid:", isValid)
}
