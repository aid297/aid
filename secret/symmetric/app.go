package symmetric

type APP struct {
	Cbc        Cbc
	Aes        Aes
	AesEncrypt AesEncrypt
	AesDecrypt AesDecrypt
	Ecb        Ecb
}
