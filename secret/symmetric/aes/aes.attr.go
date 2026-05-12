package aes

import (
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/aid297/aid/v2/secret"
)

func validateKeyBits(bits int) error {
	switch bits {
	case KeyBits128, KeyBits192, KeyBits256:
		return nil
	default:
		return fmt.Errorf("aes: invalid key bits %d, must be 128/192/256", bits)
	}
}

func toKeyBytes(bits int) int { return bits / 8 }

func KeyString(key string) secret.SymmetricAttr {
	return func(sm4Helper secret.Symmetricer) (err error) { sm4Helper.SetKeyString(key); return }
}

func KeyBytes(key []byte) secret.SymmetricAttr {
	return func(sm4Helper secret.Symmetricer) (err error) { sm4Helper.SetKeyBytes(key); return }
}

func IVString(iv string) secret.SymmetricAttr {
	return func(sm4Helper secret.Symmetricer) (err error) { sm4Helper.SetIVString(iv); return }
}

func IVBytes(iv []byte) secret.SymmetricAttr {
	return func(sm4Helper secret.Symmetricer) (err error) { sm4Helper.SetIVBytes(iv); return }
}

// KeySize 按位数生成随机 AES key（支持 128/192/256）
func KeySize(bits int, outList ...*[]byte) secret.SymmetricAttr {
	return func(symm secret.Symmetricer) error {
		if err := validateKeyBits(bits); err != nil {
			return err
		}
		helper, ok := symm.(*AESImpl)
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
			bits = KeyBits128
		}
		return RandKeyWithBits(bits, outList...)(symm)
	}
}

// RandKeyWithBits 按位数生成随机 AES key（支持 128/192/256）
func RandKeyWithBits(bits int, outList ...*[]byte) secret.SymmetricAttr {
	return func(sm4Helper secret.Symmetricer) (err error) {
		if err = validateKeyBits(bits); err != nil {
			return err
		}
		keyBytes := toKeyBytes(bits)

		key := make([]byte, keyBytes)

		if _, err = rand.Read(key); err != nil {
			return
		}

		sm4Helper.SetKeyBytes(key)
		if helper, ok := sm4Helper.(*AESImpl); ok {
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
