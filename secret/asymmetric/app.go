package asymmetric

type app struct{}

var APP app

func (*app) PEMBase64() *PEMBase64 { return NewPEMBase64() }

func (*app) RSA() *RSA { return NewRSA() }
