package secret

import (
	"github.com/aid297/aid/secret/asymmetric"
	"github.com/aid297/aid/secret/symmetric"
)

var APP struct {
	Asymmetric struct {
		RSA       asymmetric.RSA
		PEMBase64 asymmetric.PEMBase64
	}
	Symmetric struct {
		AES symmetric.AES
		CBC symmetric.CBC
		ECB symmetric.ECB
	}
}
