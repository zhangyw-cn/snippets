package signal_protocol

import (
	"bytes"
	"testing"
)

func TestX3DH(t *testing.T) {
	// Alice 和 Bob 生成密钥对
	ikA, IK_A, _ := generateKeyPair()   // Alice 的身份密钥
	ikB, IK_B, _ := generateKeyPair()   // Bob 的身份密钥
	spkB, SPK_B, _ := generateKeyPair() // Bob 的签名预密钥
	/*opkB*/ _, OPK_B, _ := generateKeyPair() // Bob 的一次性预密钥
	ekA, EK_A, _ := generateKeyPair()         // Alice 的临时密钥

	// Alice 计算共享密钥
	skAlice := x3dh(ikA, ekA, IK_B, SPK_B, OPK_B)

	// Bob 计算共享密钥
	skBob := x3dh(ikB, spkB, IK_A, EK_A, [32]byte{}) // Bob 不需要 OPK

	// 验证是否相同
	println("Session Key Match:", bytes.Equal(skAlice, skBob))
}
