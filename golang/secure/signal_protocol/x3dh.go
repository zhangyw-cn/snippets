package signal_protocol

import (
	"crypto/rand"
	"crypto/sha256"
	"io"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/hkdf"
)

// 生成 X25519 密钥对
// Curve25519 是一种高效且安全的椭圆曲线密码学（ECC）算法
// 曲线方程:   y^2=x^3 + 486662*x^2 + x 基于 Montgomery 曲线，定义在素数域F 2^255-19
// 密钥长度:   256 位（32 字节），安全性相当于 RSA-3072。
// 用途:       密钥交换（X25519）、数字签名（Ed25519）
// 性能:       比 NIST 曲线（如 P-256）更快，抗侧信道攻击。
// 前向安全性:  支持临时密钥交换（如 ECDHE），天然具备前向安全性。
func generateKeyPair() (privateKey, publicKey [32]byte, err error) {
	if _, err = io.ReadFull(rand.Reader, privateKey[:]); err != nil {
		return
	}
	curve25519.ScalarBaseMult(&publicKey, &privateKey)
	return
}

// 计算 ECDH 共享密钥
func computeDH(privateKey, publicKey [32]byte) (sharedKey [32]byte) {
	curve25519.ScalarMult(&sharedKey, &privateKey, &publicKey)
	return
}

// X3DH 密钥派生
func x3dh(
	ikA, ekA [32]byte, // Alice 的私钥
	IK_B, SPK_B, OPK_B [32]byte, // Bob 的公钥
) (sessionKey []byte) {
	// 计算 4 次 DH
	dh1 := computeDH(ikA, SPK_B)
	dh2 := computeDH(ekA, IK_B)
	dh3 := computeDH(ekA, SPK_B)
	dh4 := computeDH(ekA, OPK_B)

	// 用 HKDF 派生会话密钥
	hkdf := hkdf.New(sha256.New, append(append(dh1[:], append(dh2[:], append(dh3[:], dh4[:]...)...)...)), nil, []byte("X3DH_SS"))
	sessionKey = make([]byte, 32)
	hkdf.Read(sessionKey)
	return
}
