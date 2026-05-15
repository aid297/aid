package sm4

import (
	"crypto/rand"

	"github.com/aid297/aid/v2/secret"
)

func KeyString(key string) secret.SymmetricAttr {
	return func(symmetricer secret.Symmetric) (err error) { symmetricer.SetKeyString(key); return }
}

func KeyBytes(key []byte) secret.SymmetricAttr {
	return func(symmetricer secret.Symmetric) (err error) { symmetricer.SetKeyBytes(key); return }
}

func IVString(iv string) secret.SymmetricAttr {
	return func(symmetricer secret.Symmetric) (err error) { symmetricer.SetIVString(iv); return }
}

func IVBytes(iv []byte) secret.SymmetricAttr {
	return func(symmetricer secret.Symmetric) (err error) { symmetricer.SetIVBytes(iv); return }
}

func RandKey(outList ...*[]byte) secret.SymmetricAttr {
	return func(symmetricer secret.Symmetric) (err error) {
		key := make([]byte, 16)

		if _, err = rand.Read(key); err != nil {
			return
		}

		symmetricer.SetKeyBytes(key)

		if len(outList) > 0 && outList[0] != nil {
			*outList[0] = make([]byte, 16)
			copy(*outList[0], key)
		}

		return
	}
}

func RandIV(outList ...*[]byte) secret.SymmetricAttr {
	return func(symmetricer secret.Symmetric) (err error) {
		iv := make([]byte, 16)

		if _, err = rand.Read(iv); err != nil {
			return
		}

		symmetricer.SetIVBytes(iv)

		if len(outList) > 0 && outList[0] != nil {
			*outList[0] = make([]byte, 16)
			copy(*outList[0], iv)
		}

		return
	}
}

func AlgorithmECB() secret.SymmetricAttr {
	return func(symmetricer secret.Symmetric) (err error) { return symmetricer.SetAlgorithm("ECB") }
}

func AlgorithmCBC() secret.SymmetricAttr {
	return func(symmetricer secret.Symmetric) (err error) { return symmetricer.SetAlgorithm("CBC") }
}

func AlgorithmCTR() secret.SymmetricAttr {
	return func(symmetricer secret.Symmetric) (err error) { return symmetricer.SetAlgorithm("CTR") }
}

func AlgorithmGCM() secret.SymmetricAttr {
	return func(symmetricer secret.Symmetric) (err error) { return symmetricer.SetAlgorithm("GCM") }
}

// RandCTRNonce 生成随机 nonce（适用于 CTR 模式：16字节）
func RandCTRNonce(outList ...*[]byte) secret.SymmetricAttr {
	return func(symmetricer secret.Symmetric) (err error) {
		nonce := make([]byte, 16)

		if _, err = rand.Read(nonce); err != nil {
			return
		}

		symmetricer.SetIVBytes(nonce)

		if len(outList) > 0 && outList[0] != nil {
			*outList[0] = make([]byte, 16)
			copy(*outList[0], nonce)
		}

		return
	}
}

// RandGCMNonce 生成 12 字节随机 nonce（适用于 GCM 模式）
func RandGCMNonce(outList ...*[]byte) secret.SymmetricAttr {
	return func(symmetricer secret.Symmetric) (err error) {
		nonce := make([]byte, 12)

		if _, err = rand.Read(nonce); err != nil {
			return
		}

		symmetricer.SetIVBytes(nonce)

		if len(outList) > 0 && outList[0] != nil {
			*outList[0] = make([]byte, 12)
			copy(*outList[0], nonce)
		}

		return
	}
}
