package symmetric

type app struct{}

var APP app

func (*app) AES(sail string) *AES { return NewAES(sail) }

func (*app) CBC() *CBC { return NewCBC() }

func (*app) ECB() *ECB { return NewECB() }
