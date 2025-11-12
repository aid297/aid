package symmetric

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/aid297/aid/compression"
)

func TestEcbDemo(t *testing.T) {
	t.Run("ecb加密", func(t *testing.T) {
		var (
			err                                        error
			zipped, encrypted, decrypted, decompressed []byte
			ciphertext                                 string
		)

		a := NewAes("tjp5OPIU1ETF5s33fsLWdA==")
		aesEncrypt := a.NewEncrypt().GetEncrypt()
		openKey := aesEncrypt.GetOpenKey()
		aesDecrypt := a.NewDecrypt(openKey).GetDecrypt()

		if plaintext, jsonErr := json.Marshal([]map[string]any{{"name": "张三", "age": 18}, {"name": "李四", "age": 20}}); jsonErr != nil {
			t.Errorf("[ECB] json marshal: %v", jsonErr)
		} else {
			log.Printf("[ECB] json marshal: %s\n", plaintext)

			// encrypt
			// encrypt step1: zip
			zipped, err = compression.NewZlib().Compress(plaintext)
			if err != nil {
				t.Errorf("[ECB] compressing: %v\n", err)
			}

			// encrypt step2: aes-ecb-encrypt
			encrypted, err = Ecb{}.Encrypt(aesEncrypt.GetAesKey(), zipped)
			if err != nil {
				t.Errorf("[ECB] encrypting: %v", err)
			}

			// encrypt step3: encode base64
			ciphertext = base64.StdEncoding.EncodeToString(encrypted)
			fmt.Printf("encrypted: %s", ciphertext)

			// decrypt step1: decrypt
			decrypted, err = Ecb{}.Decrypt(aesDecrypt.GetAesKey(), encrypted)
			if err != nil {
				t.Errorf("[ECB] decrypting: %v", err)
			}

			// decrypt step2: decompress
			decompressed, err = compression.NewZlib().Decompress(decrypted)
			if err != nil {
				t.Errorf("[ECB] decompressing: %v", err)
			}
			fmt.Printf("decrypted: %s\n", decompressed)
		}
	})
}
