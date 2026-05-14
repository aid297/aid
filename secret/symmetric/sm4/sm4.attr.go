package sm4

import (
	"crypto/rand"

	"github.com/aid297/aid/v2/secret"
)

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

func RandKey(outList ...*[]byte) secret.SymmetricAttr {
	return func(symmetricer secret.Symmetricer) (err error) {
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
	return func(symmetricer secret.Symmetricer) (err error) {
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
