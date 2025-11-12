package secret

import (
	"github.com/aid297/aid/secret/asymmetric"
	"github.com/aid297/aid/secret/symmetric"
)

type APP Launch

type Launch struct {
	Asymmetric struct {
		Rsa       asymmetric.Rsa
		PemBase64 asymmetric.PemBase64
	}
	Symmetric struct {
		Cbc        symmetric.Cbc
		Aes        symmetric.Aes
		AesEncrypt symmetric.AesEncrypt
		AesDecrypt symmetric.AesDecrypt
		Ecb        symmetric.Ecb
	}
}
