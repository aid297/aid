package aes

import (
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/aid297/aid/v2/secret"
)

func validateKeyBits(bits int) error {
	switch bits {
	case AES128, AES192, AES256:
		return nil
	default:
		return fmt.Errorf("aes: invalid key bits %d, must be 128/192/256", bits)
	}
}

func toKeyBytes(bits int) int { return bits / 8 }

func KeyString(key string) secret.SymmetricAttr {
	return func(symmetricer secret.Symmetricer) (err error) { symmetricer.SetKeyString(key); return }
}

func KeyBytes(key []byte) secret.SymmetricAttr {
	return func(symmetricer secret.Symmetricer) (err error) { symmetricer.SetKeyBytes(key); return }
}

func IVString(iv string) secret.SymmetricAttr {
	return func(symmetricer secret.Symmetricer) (err error) { symmetricer.SetIVString(iv); return }
}

func IVBytes(iv []byte) secret.SymmetricAttr {
	return func(symmetricer secret.Symmetricer) (err error) { symmetricer.SetIVBytes(iv); return }
}

// KeySize 按位数生成随机 AES key（支持 128/192/256）
func KeySize(bits int, outList ...*[]byte) secret.SymmetricAttr {
	return func(symmetricer secret.Symmetricer) error {
		if err := validateKeyBits(bits); err != nil {
			return err
		}
		helper, ok := symmetricer.(*AESImpl)
		if !ok {
			return errors.New("aes: KeySize only supports AESImpl")
		}
		helper.keyBits = bits
		if len(outList) > 0 && outList[0] != nil {
			*outList[0] = make([]byte, toKeyBytes(bits))
		}
		return nil
	}
}

func RandKey(outList ...*[]byte) secret.SymmetricAttr {
	return func(symm secret.Symmetricer) (err error) {
		helper, ok := symm.(*AESImpl)
		if !ok {
			return errors.New("aes: RandKey only supports AESImpl")
		}
		bits := helper.keyBits
		if bits == 0 {
			bits = AES128
		}
		return RandKeyWithBits(bits, outList...)(symm)
	}
}

// RandKeyWithBits 按位数生成随机 AES key（支持 128/192/256）
func RandKeyWithBits(bits int, outList ...*[]byte) secret.SymmetricAttr {
	return func(symmetricer secret.Symmetricer) (err error) {
		if err = validateKeyBits(bits); err != nil {
			return err
		}
		keyBytes := toKeyBytes(bits)

		key := make([]byte, keyBytes)

		if _, err = rand.Read(key); err != nil {
			return
		}

		symmetricer.SetKeyBytes(key)
		if helper, ok := symmetricer.(*AESImpl); ok {
			helper.keyBits = bits
		}

		if len(outList) > 0 && outList[0] != nil {
			*outList[0] = make([]byte, keyBytes)
			copy(*outList[0], key)
		}

		return
	}
}

func RandIV(outList ...*[]byte) secret.SymmetricAttr {
	return func(sm4Helper secret.Symmetricer) (err error) {
		iv := make([]byte, 16)

		if _, err = rand.Read(iv); err != nil {
			return
		}

		sm4Helper.SetIVBytes(iv)

		if len(outList) > 0 && outList[0] != nil {
			*outList[0] = make([]byte, 16)
			copy(*outList[0], iv)
		}

		return
	}
}

func AlgorithmECB() secret.SymmetricAttr {
	return func(symmetricer secret.Symmetricer) (err error) { return symmetricer.SetAlgorithm("ECB") }
}

func AlgorithmCBC() secret.SymmetricAttr {
	return func(symmetricer secret.Symmetricer) (err error) { return symmetricer.SetAlgorithm("CBC") }
}

func AlgorithmCTR() secret.SymmetricAttr {
	return func(symmetricer secret.Symmetricer) (err error) { return symmetricer.SetAlgorithm("CTR") }
}

func AlgorithmGCM() secret.SymmetricAttr {
	return func(symmetricer secret.Symmetricer) (err error) { return symmetricer.SetAlgorithm("GCM") }
}

// RandCTRNonce 生成随机 nonce（适用于 CTR 模式：16字节）
func RandCTRNonce(outList ...*[]byte) secret.SymmetricAttr {
	return func(symmetricer secret.Symmetricer) (err error) {
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
	return func(symmetricer secret.Symmetricer) (err error) {
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
