## *Secret* 加密模块

1. 非对称加密（*SM2*）

   1. 生成种子
      ```go
      func TestGenerateKeyPair(t *testing.T) {
      	var (
      		err                        error
      		sem                        secret.Semener
      		pubKeyBase64, priKeyBase64 string
      		pubKeyBytes, priKeyBytes   []byte
      	)
      
      	if sem, err = sm2.NewSem(); err != nil {
      		t.Fatalf("生成种子失败：%v", err)
      	}
      
      	if err = sem.GeneratePriKey(); err != nil {
      		t.Fatalf("生成私钥失败：%v", err)
      	}
      
      	if pubKeyBase64, err = sem.GetPubKeyBase64(); pubKeyBase64 == "" {
      		t.Fatalf("公钥长度为空：base64")
      	}
      
      	if pubKeyBytes, err = sem.GetPubKeyBytes(); len(pubKeyBytes) == 0 {
      		t.Fatalf("公钥长度为空：bytes")
      	}
      
      	if priKeyBase64, err = sem.GetPriKeyBase64(); priKeyBase64 == "" {
      		t.Fatalf("私钥长度为空：base64")
      	}
      
      	if priKeyBytes, err = sem.GetPriKeyBytes(); len(priKeyBytes) == 0 {
      		t.Fatalf("私钥长度为空：bytes")
      	}
      
      	t.Logf("公钥：%s\n", pubKeyBase64)
      	t.Logf("私钥：%s\n", priKeyBase64)
      }
      ```

      **说明：**sm2.NewSem()这个方法其实已经验证了公钥和密钥，如果NewSem()这个方法没有带有私钥参数，则会自动生成随机

   2. 加解密
      ```go
      func TestEncryptDecrypt(t *testing.T) {
      	var (
      		err                error
      		semA, semB         secret.Semener
      		sm2A, sm2B         secret.Asymmetricer
      		sm2SemAPriKeyBytes []byte
      		plainText          = []byte("hello, SM2 非对称加密测试!")
      		cipherBase64       string
      		decrypted          []byte
      	)
      
      	if semA, err = sm2.NewSem(); err != nil {
      		t.Fatalf("生成加密种子失败：%v", err)
      	}
      
      	sm2A = sm2.New(semA)
      
      	if cipherBase64, err = sm2A.Encrypt(plainText); err != nil {
      		t.Fatalf("加密错误：%v", err)
      	}
      	t.Logf("加密结果：%s\n", cipherBase64)
      
      	if sm2SemAPriKeyBytes, err = semA.GetPriKeyBytes(); err != nil {
      		t.Fatalf("获取私钥失败：%v", err)
      	}
      
      	if semB, err = sm2.NewSem(sm2.PriKeyBytes(sm2SemAPriKeyBytes)); err != nil {
      		t.Fatalf("生成解密种子失败：%v", err)
      	}
      	a, _ := semB.GetPriKeyBase64()
      	b, _ := semB.GetPubKeyBase64()
      	t.Logf("解密种子私钥：%s, 公钥：%s", a, b)
      
      	sm2B = sm2.New(semB)
      
      	if decrypted, err = sm2B.Decrypt(cipherBase64); err != nil {
      		t.Fatalf("解密错误：%v", err)
      	}
      	t.Logf("解密结果：%s", decrypted)
      
      	if !bytes.Equal(plainText, decrypted) {
      		t.Fatalf("比对结果不匹配")
      	}
      }
      ```

      **说明：**NewSem(PriKeyBytes([]byte(...)))手动设置私钥并

   3. 签名&验证
      ```go
      func TestSignVerify(t *testing.T) {
      	var (
      		err       error
      		semA      secret.Semener
      		sm2Helper secret.Asymmetricer
      		data      = []byte("hello, SM2 数字签名测试!")
      		sigHex    string
      		ok        bool
      	)
      
      	if semA, err = sm2.NewSem(); err != nil {
      		t.Fatalf("生成种子失败：%v", err)
      	}
      
      	sm2Helper = sm2.New(semA)
      
      	if sigHex, err = sm2Helper.Sign(data); err != nil {
      		t.Fatalf("签名失败：%v", err)
      	}
      	t.Logf("签名内容(hex): %s", sigHex)
      
      	if ok, err = sm2Helper.Verify(data, sigHex); err != nil {
      		t.Fatalf("验证失败：%v", err)
      	}
      	if !ok {
      		t.Fatal("签名验证失败")
      	}
      
      	t.Logf("验证成功")
      }
      
      func TestVerifyWithWrongData(t *testing.T) {
      	var (
      		err        error
      		semA, semB secret.Semener
      		sm2A, sm2B secret.Asymmetricer
      		data       = []byte("original data")
      		sigHex     string
      		ok         bool
      	)
      	if semA, err = sm2.NewSem(); err != nil {
      		t.Fatalf("生成种子失败：%v", err)
      	}
      
      	if err = sm2.MustGeneratePriKey(semA); err != nil {
      		t.Fatalf("%v", err)
      	}
      
      	sm2A = sm2.New(semA)
      
      	if sigHex, err = sm2A.Sign(data); err != nil {
      		t.Fatalf("签名失败：%v", err)
      	}
      
      	semB, err = sm2.NewSem(sm2.PriKey(semA.GetPriKey()))
      	if err != nil {
      		t.Fatalf("生成解密种子失败：%v", err)
      	}
      
      	sm2B = sm2.New(semB)
      
      	if ok, err = sm2B.Verify([]byte("tampered data"), sigHex); err != nil {
      		t.Fatalf("验证失败：%v", err)
      	}
      	if ok {
      		t.Fatal("篡改数据的验证应失败")
      	}
      	t.Logf("篡改数据正确被拒绝")
      }
      ```

2. 非对称加密（*RSA*）

   1. 生成种子
      ```go
      func TestGenerateKeyPair(t *testing.T) {
      	var (
      		err                        error
      		sem                        secret.Semener
      		pubKeyBase64, priKeyBase64 string
      		pubKeyBytes, priKeyBytes   []byte
      	)
      
      	if sem, err = rsa.NewSem(); err != nil {
      		t.Fatalf("生成种子(RSA)失败：%v", err)
      	}
      
      	if err = sem.GeneratePriKey(); err != nil {
      		t.Fatalf("生成私钥失败：%v", err)
      	}
      
      	if pubKeyBase64, err = sem.GetPubKeyBase64(); pubKeyBase64 == "" {
      		t.Fatalf("公钥长度为空：base64")
      	}
      
      	if pubKeyBytes, err = sem.GetPubKeyBytes(); len(pubKeyBytes) == 0 {
      		t.Fatalf("公钥长度为空：bytes")
      	}
      
      	if priKeyBase64, err = sem.GetPriKeyBase64(); priKeyBase64 == "" {
      		t.Fatalf("私钥长度为空：base64")
      	}
      
      	if priKeyBytes, err = sem.GetPriKeyBytes(); len(priKeyBytes) == 0 {
      		t.Fatalf("私钥长度为空：bytes")
      	}
      
      	t.Logf("公钥：%s\n", pubKeyBase64)
      	t.Logf("私钥：%s\n", priKeyBase64)
      }
      ```
      
   2. 加解密
      ```go
      func TestEncryptDecrypt(t *testing.T) {
      	var (
      		err                error
      		semA, semB         secret.Semener
      		rsaA, rsaB         secret.Asymmetricer
      		rsaSemAPriKeyBytes []byte
      		plainText          = []byte("hello, RSA 非对称加密测试!")
      		cipherBase64       string
      		decrypted          []byte
      	)
      
      	if semA, err = rsa.NewSem(); err != nil {
      		t.Fatalf("生成加密种子失败：%v", err)
      	}
      	rsaA = rsa.New(semA)
      
      	if cipherBase64, err = rsaA.Encrypt(plainText); err != nil {
      		t.Fatalf("加密错误：%v", err)
      	}
      	t.Logf("加密结果：%s\n", cipherBase64)
      
      	if rsaSemAPriKeyBytes, err = semA.GetPriKeyBytes(); err != nil {
      		t.Fatalf("获取私钥失败：%v", err)
      	}
      
      	if semB, err = rsa.NewSem(rsa.PriKeyBytes(rsaSemAPriKeyBytes)); err != nil {
      		t.Fatalf("生成解密种子失败：%v", err)
      	}
      	a, _ := semB.GetPriKeyBase64()
      	b, _ := semB.GetPubKeyBase64()
      	t.Logf("解密种子私钥：%s, 公钥：%s", a, b)
      
      	rsaB = rsa.New(semB)
      
      	if decrypted, err = rsaB.Decrypt(cipherBase64); err != nil {
      		t.Fatalf("解密错误：%v", err)
      	}
      	t.Logf("解密结果：%s", decrypted)
      
      	if !bytes.Equal(plainText, decrypted) {
      		t.Fatalf("比对结果不匹配")
      	}
      }
      ```
      
   2. 签名&验证
      ```go
      func TestSignVerify(t *testing.T) {
      	var (
      		err       error
      		semA      secret.Semener
      		rsaHelper secret.Asymmetricer
      		data      = []byte("hello, RSA 数字签名测试!")
      		sigHex    string
      		ok        bool
      	)
      
      	if semA, err = rsa.NewSem(); err != nil {
      		t.Fatalf("生成种子失败：%v", err)
      	}
      
      	rsaHelper = rsa.New(semA)
      
      	if sigHex, err = rsaHelper.Sign(data); err != nil {
      		t.Fatalf("签名失败：%v", err)
      	}
      	t.Logf("签名内容(hex): %s", sigHex)
      
      	if ok, err = rsaHelper.Verify(data, sigHex); err != nil {
      		t.Fatalf("验证失败：%v", err)
      	}
      	if !ok {
      		t.Fatal("签名验证失败")
      	}
      
      	t.Logf("验证成功")
      }
      
      func TestVerifyWithWrongData(t *testing.T) {
      	var (
      		err        error
      		semA, semB secret.Semener
      		rsaA, rsaB secret.Asymmetricer
      		data       = []byte("original data")
      		sigHex     string
      		ok         bool
      	)
      	if semA, err = rsa.NewSem(); err != nil {
      		t.Fatalf("生成种子失败：%v", err)
      	}
      
      	if err = rsa.MustGeneratePriKey(semA); err != nil {
      		t.Fatalf("%v", err)
      	}
      
      	rsaA = rsa.New(semA)
      	if sigHex, err = rsaA.Sign(data); err != nil {
      		t.Fatalf("签名失败：%v", err)
      	}
      
      	semB, err = rsa.NewSem(rsa.PriKey(semA.GetPriKey()))
      	if err != nil {
      		t.Fatalf("生成解密种子失败：%v", err)
      	}
      
      	rsaB = rsa.New(semB)
      	if ok, err = rsaB.Verify([]byte("tampered data"), sigHex); err != nil {
      		t.Fatalf("验证失败：%v", err)
      	}
      	if ok {
      		t.Fatal("篡改数据的验证应失败")
      	}
      	t.Logf("篡改数据正确被拒绝")
      }
      ```

3. 对称加密（*SM4-ECB*）

   1. 基本用法
      ```go
      func TestECB(t *testing.T) {
      	sm4Helper, err := sm4.New(sm4.KeyBytes(testKey))
      	if err != nil {
      		t.Fatalf("创建 ECB 对象失败：%v", err)
      	}
      	cipherText, err := sm4Helper.EncryptECB(testPlain)
      	if err != nil {
      		t.Fatalf("EncryptECB failed: %v", err)
      	}
      	plain, err := sm4Helper.DecryptECB(cipherText)
      	if err != nil {
      		t.Fatalf("DecryptECB failed: %v", err)
      	}
      	if !bytes.Equal(plain, testPlain) {
      		t.Fatalf("ECB mismatch: got %s, want %s", plain, testPlain)
      	}
      	t.Logf("ECB OK: %s", plain)
      }
      
      func TestECBBase64(t *testing.T) {
      	sm4Helper, err := sm4.New(sm4.KeyBytes(testKey))
      	if err != nil {
      		t.Fatalf("创建 ECB 对象失败：%v", err)
      	}
      	b64, err := sm4Helper.EncryptECBBase64(testPlain)
      	if err != nil {
      		t.Fatalf("EncryptECBBase64 failed: %v", err)
      	}
      	t.Logf("ECB base64: %s", b64)
      	plain, err := sm4Helper.DecryptECBBase64(b64)
      	if err != nil {
      		t.Fatalf("DecryptECBBase64 failed: %v", err)
      	}
      	if !bytes.Equal(plain, testPlain) {
      		t.Fatalf("ECB base64 mismatch: got %s, want %s", plain, testPlain)
      	}
      	t.Logf("ECB base64 OK: %s", plain)
      }
      ```

4. 对称加密（*SM4-CBC*）

   1. 基本用法
      ```go
      func TestCBC(t *testing.T) {
      	sm4Helper, err := sm4.New(sm4.KeyBytes(testKey), sm4.IVBytes(testIV))
      	if err != nil {
      		t.Fatalf("创建 CBC 对象失败：%v", err)
      	}
      	cipherText, err := sm4Helper.EncryptCBC(testPlain)
      	if err != nil {
      		t.Fatalf("EncryptCBC failed: %v", err)
      	}
      	plain, err := sm4Helper.DecryptCBC(cipherText)
      	if err != nil {
      		t.Fatalf("DecryptCBC failed: %v", err)
      	}
      	if !bytes.Equal(plain, testPlain) {
      		t.Fatalf("CBC mismatch: got %s, want %s", plain, testPlain)
      	}
      	t.Logf("CBC OK: %s", plain)
      }
      
      func TestCBCBase64(t *testing.T) {
      	sm4Helper, err := sm4.New(sm4.KeyBytes(testKey), sm4.IVBytes(testIV))
      	if err != nil {
      		t.Fatalf("创建 CBC 对象失败：%v", err)
      	}
      	b64, err := sm4Helper.EncryptCBCBase64(testPlain)
      	if err != nil {
      		t.Fatalf("EncryptCBCBase64 failed: %v", err)
      	}
      	t.Logf("CBC base64: %s", b64)
      	plain, err := sm4Helper.DecryptCBCBase64(b64)
      	if err != nil {
      		t.Fatalf("DecryptCBCBase64 failed: %v", err)
      	}
      	if !bytes.Equal(plain, testPlain) {
      		t.Fatalf("CBC base64 mismatch: got %s, want %s", plain, testPlain)
      	}
      	t.Logf("CBC base64 OK: %s", plain)
      }
      
      func TestInvalidKey(t *testing.T) {
      	sm4Helper, err := sm4.New(sm4.KeyString("shortkey"))
      	if err != nil {
      		t.Fatalf("创建对象失败：%v", err)
      	}
      	_, err = sm4Helper.EncryptECB(testPlain)
      	if err == nil {
      		t.Fatal("expected error for invalid key length")
      	}
      	t.Logf("Invalid key error: %v", err)
      }
      ```

5. 对称加密（*AES128-ECB*）

   1. 基本用法

      ```go
      func TestECB(t *testing.T) {
      	var (
      		testKey   = []byte("1234567890abcdef")
      		testPlain = []byte("hello, AES encrypt test!")
      	)
      
      	aesHelper, err := aes.New(aes.KeyBytes(testKey))
      	if err != nil {
      		t.Fatalf("创建 ECB 对象失败：%v", err)
      	}
      
      	cipherText, err := aesHelper.EncryptECB(testPlain)
      	if err != nil {
      		t.Fatalf("EncryptECB failed: %v", err)
      	}
      
      	plain, err := aesHelper.DecryptECB(cipherText)
      	if err != nil {
      		t.Fatalf("DecryptECB failed: %v", err)
      	}
      
      	if !bytes.Equal(plain, testPlain) {
      		t.Fatalf("ECB mismatch: got %s, want %s", plain, testPlain)
      	}
      }
      
      func TestECBBase64(t *testing.T) {
      	var (
      		testKey   = []byte("1234567890abcdef")
      		testPlain = []byte("hello, AES encrypt test!")
      	)
      
      	aesHelper, err := aes.New(aes.KeyBytes(testKey))
      	if err != nil {
      		t.Fatalf("创建 ECB 对象失败：%v", err)
      	}
      
      	b64, err := aesHelper.EncryptECBBase64(testPlain)
      	if err != nil {
      		t.Fatalf("EncryptECBBase64 failed: %v", err)
      	}
      
      	plain, err := aesHelper.DecryptECBBase64(b64)
      	if err != nil {
      		t.Fatalf("DecryptECBBase64 failed: %v", err)
      	}
      
      	if !bytes.Equal(plain, testPlain) {
      		t.Fatalf("ECB base64 mismatch: got %s, want %s", plain, testPlain)
      	}
      }
      ```

6. 对称加密（*AES128-CBC*）

   1. 基本用法
      ```go
      func TestCBC(t *testing.T) {
      	var (
      		testKey   = []byte("1234567890abcdef")
      		testIV    = []byte("abcdef1234567890")
      		testPlain = []byte("hello, AES encrypt test!")
      	)
      
      	aesHelper, err := aes.New(aes.KeyBytes(testKey), aes.IVBytes(testIV))
      	if err != nil {
      		t.Fatalf("创建 CBC 对象失败：%v", err)
      	}
      
      	cipherText, err := aesHelper.EncryptCBC(testPlain)
      	if err != nil {
      		t.Fatalf("EncryptCBC failed: %v", err)
      	}
      
      	plain, err := aesHelper.DecryptCBC(cipherText)
      	if err != nil {
      		t.Fatalf("DecryptCBC failed: %v", err)
      	}
      
      	if !bytes.Equal(plain, testPlain) {
      		t.Fatalf("CBC mismatch: got %s, want %s", plain, testPlain)
      	}
      }
      
      func TestCBCBase64(t *testing.T) {
      	var (
      		testKey   = []byte("1234567890abcdef")
      		testIV    = []byte("abcdef1234567890")
      		testPlain = []byte("hello, AES encrypt test!")
      	)
      
      	aesHelper, err := aes.New(aes.KeyBytes(testKey), aes.IVBytes(testIV))
      	if err != nil {
      		t.Fatalf("创建 CBC 对象失败：%v", err)
      	}
      
      	b64, err := aesHelper.EncryptCBCBase64(testPlain)
      	if err != nil {
      		t.Fatalf("EncryptCBCBase64 failed: %v", err)
      	}
      
      	plain, err := aesHelper.DecryptCBCBase64(b64)
      	if err != nil {
      		t.Fatalf("DecryptCBCBase64 failed: %v", err)
      	}
      
      	if !bytes.Equal(plain, testPlain) {
      		t.Fatalf("CBC base64 mismatch: got %s, want %s", plain, testPlain)
      	}
      }
      
      func TestInvalidKey(t *testing.T) {
      	var testPlain = []byte("hello, AES encrypt test!")
      
      	aesHelper, err := aes.New(aes.KeyString("shortkey"))
      	if err != nil {
      		t.Fatalf("创建对象失败：%v", err)
      	}
      
      	_, err = aesHelper.EncryptECB(testPlain)
      	if err == nil {
      		t.Fatal("expected error for invalid key length")
      	}
      }
      ```
      
   2. *AES192*和*AES256*支持

      ```go
      func TestCBC192And256(t *testing.T) {
      	var (
      		testKey24 = []byte("1234567890abcdefghijklmn")
      		testKey32 = []byte("1234567890abcdefghijklmnopqrstuv")
      		testIV    = []byte("abcdef1234567890")
      		testPlain = []byte("hello, AES encrypt test!")
      		tests     = []struct {
      			name string
      			key  []byte
      		}{
      			{name: "AES-192", key: testKey24},
      			{name: "AES-256", key: testKey32},
      		}
      	)
      
      	for _, tt := range tests {
      		t.Run(tt.name, func(t *testing.T) {
      			aesHelper, err := aes.New(aes.KeyBytes(tt.key), aes.IVBytes(testIV))
      			if err != nil {
      				t.Fatalf("创建 CBC 对象失败：%v", err)
      			}
      			cipherText, err := aesHelper.EncryptCBC(testPlain)
      			if err != nil {
      				t.Fatalf("EncryptCBC failed: %v", err)
      			}
      			plain, err := aesHelper.DecryptCBC(cipherText)
      			if err != nil {
      				t.Fatalf("DecryptCBC failed: %v", err)
      			}
      			if !bytes.Equal(plain, testPlain) {
      				t.Fatalf("CBC mismatch: got %s, want %s", plain, testPlain)
      			}
      		})
      	}
      }
      
      func TestRandKeyWithBits(t *testing.T) {
      	var (
      		key       []byte
      		testPlain = []byte("hello, AES encrypt test!")
      	)
      
      	aesHelper, err := aes.New(aes.RandKeyWithBits(aes.KeyBits256, &key), aes.RandIV())
      	if err != nil {
      		t.Fatalf("创建对象失败：%v", err)
      	}
      
      	if len(key) != 32 {
      		t.Fatalf("随机 key 长度错误: got %d, want 32", len(key))
      	}
      
      	if _, err = aesHelper.EncryptCBC(testPlain); err != nil {
      		t.Fatalf("EncryptCBC failed: %v", err)
      	}
      }
      
      func TestRandKeyWithKeySize(t *testing.T) {
      	var (
      		key       []byte
      		testPlain = []byte("hello, AES encrypt test!")
      	)
      
      	aesHelper, err := aes.New(aes.KeySize(aes.KeyBits192), aes.RandKey(&key), aes.RandIV())
      	if err != nil {
      		t.Fatalf("创建对象失败：%v", err)
      	}
      
      	if len(key) != 24 {
      		t.Fatalf("KeySize+RandKey 长度错误: got %d, want 24", len(key))
      	}
      
      	if _, err = aesHelper.EncryptCBC(testPlain); err != nil {
      		t.Fatalf("EncryptCBC failed: %v", err)
      	}
      }
      
      func TestRandKeyWithBitsInvalid(t *testing.T) {
      	_, err := aes.New(aes.RandKeyWithBits(111), aes.RandIV())
      	if err == nil {
      		t.Fatal("expected error for invalid key bits")
      	}
      }
      ```

7. 组合用法

   1. *SM2+SM4-CBC*加密文件和大文件加密

      ```go
      package main
      
      import (
      	"crypto/rand"
      	"errors"
      	"io"
      	"log"
      	"os"
      
      	"github.com/aid297/aid/v2/secret"
      	"github.com/aid297/aid/v2/secret/asymmetric/sm2"
      	"github.com/aid297/aid/v2/secret/symmetric/sm4"
      )
      
      // TestFileEncrypt 1. 文件加密/解密演示
      func TestFileEncrypt() {
      	// 生成密钥对
      	var (
      		err                        error
      		semEncrypter, semDecrypter secret.Semener
      		sm2Encrypter, sm2Decrypter secret.Asymmetricer
      		semPriKeyEncrypter         []byte
      		plainFile                  = "/tmp/sm2_test_plain.txt"
      		encryptedFile              = "/tmp/sm2_test_encrypted.bin"
      		decryptedFile              = "/tmp/sm2_test_decrypted.txt"
      		sm4Encrypter, sm4Decrypter secret.Symmetricer
      		key                        []byte
      		iv                         []byte
      	)
      
      	// 生成加密种子
      	if semEncrypter, err = sm2.NewSem(); err != nil {
      		log.Fatalf("生成加密种(SM2)子失败：%v", err)
      	}
      	sm2Encrypter = sm2.New(semEncrypter) // 生成加密器(SM2)
      
      	// 准备测试文件
      	if err = os.WriteFile(plainFile, []byte("这是需要加密的文件内容，可以是任意大小的数据。"), 0644); err != nil {
      		log.Fatalf("写入测试文件失败：%v", err)
      	}
      
      	// 加密
      	if sm4Encrypter, err = sm4.New(sm4.RandKey(&key), sm4.RandIV(&iv)); err != nil {
      		log.Fatalf("生成加密器(SM4)失败：%v", err)
      	}
      	if err = sm4Encrypter.EncryptCBCFile(plainFile, encryptedFile, sm2Encrypter); err != nil {
      		log.Fatalf("加密失败：%v", err)
      	}
      	log.Printf("加密成功：%s", encryptedFile)
      
      	// 解密
      	if semPriKeyEncrypter, err = semEncrypter.GetPriKeyBytes(); err != nil {
      		log.Fatalf("获取加密种子(SM2)失败：%v", err)
      	}
      	// 通过加密种子(SM2)制作解密种子
      	if semDecrypter, err = sm2.NewSem(sm2.PriKeyBytes(semPriKeyEncrypter)); err != nil {
      		log.Fatalf("生成解密种子(SM2)失败：%v", err)
      	}
      	sm2Decrypter = sm2.New(semDecrypter)
      	if sm4Decrypter, err = sm4.New(sm4.KeyBytes(key), sm4.IVBytes(iv)); err != nil {
      		log.Fatalf("生成解密器(SM4)失败：%v", err)
      	}
      	if err = sm4Decrypter.DecryptCBCFile(encryptedFile, decryptedFile, sm2Decrypter); err != nil {
      		log.Fatalf("解密文件失败：%v", err)
      	}
      	log.Printf("解密成功：%s", decryptedFile)
      
      	// 对比结果
      	original, _ := os.ReadFile(plainFile)
      	decrypted, _ := os.ReadFile(decryptedFile)
      	if string(original) == string(decrypted) {
      		log.Print("✓ 文件内容一致，加解密正确")
      	} else {
      		log.Fatal("✗ 文件内容不一致")
      	}
      
      	cleanupTestFiles(plainFile, encryptedFile, decryptedFile)
      }
      
      // TestLargeFileEncrypt 2. 大文件加密/解密演示（流式）
      func TestLargeFileEncrypt() {
      	var (
      		err                        error
      		semEncrypter, semDecrypter secret.Semener
      		sm2Encrypter, sm2Decrypter secret.Asymmetricer
      		sm4Encrypter, sm4Decrypter secret.Symmetricer
      		semEncrypterPriKeyBytes    []byte
      		plainFile                        = "/tmp/sm2_test_large_plain.bin"
      		encryptedFile                    = "/tmp/sm2_test_large_encrypted.bin"
      		decryptedFile                    = "/tmp/sm2_test_large_decrypted.bin"
      		fileSize                   int64 = 64 * 1024 * 1024 // 64MB 演示
      		key, iv                    []byte
      	)
      
      	if semEncrypter, err = sm2.NewSem(); err != nil {
      		log.Fatalf("生成种子失败：%v", err)
      	}
      	sm2Encrypter = sm2.New(semEncrypter)
      
      	// 1) 生成大文件（随机数据）
      	func() {
      		var (
      			f   *os.File
      			buf = make([]byte, 1024*1024)
      			wr  int64
      		)
      		if f, err = os.Create(plainFile); err != nil {
      			return
      		}
      		defer func() { _ = f.Close() }()
      
      		for wr < fileSize {
      			need := min(fileSize-wr, int64(len(buf)))
      			if _, err = rand.Read(buf[:need]); err != nil {
      				return
      			}
      			if _, err = f.Write(buf[:need]); err != nil {
      				return
      			}
      			wr += need
      		}
      	}()
      	if err != nil {
      		log.Fatalf("生成大文件失败：%v", err)
      	}
      
      	// 2) 流式加密
      	if sm4Encrypter, err = sm4.New(sm4.RandKey(&key), sm4.RandIV(&iv)); err != nil {
      		log.Fatalf("生成加密器(SM4)失败：%v", err)
      	}
      	if err = sm4Encrypter.EncryptCBCLargeFile(plainFile, encryptedFile, sm2Encrypter); err != nil {
      		log.Fatalf("大文件加密失败：%v", err)
      	}
      	log.Printf("大文件加密成功：%s", encryptedFile)
      
      	// 3) 流式解密
      	if sm4Decrypter, err = sm4.New(sm4.KeyBytes(key), sm4.IVBytes(iv)); err != nil {
      		log.Fatalf("生成解密器(SM4)失败：%v", err)
      	}
      
      	// 通过加密种子获取解密种子所需的私钥字节（此时两个种子应该一致）
      	if semEncrypterPriKeyBytes, err = semEncrypter.GetPriKeyBytes(); err != nil {
      		log.Fatalf("获取加密种子密钥失败：%v", err)
      	}
      	if semDecrypter, err = sm2.NewSem(sm2.PriKeyBytes(semEncrypterPriKeyBytes)); err != nil {
      		log.Fatalf("生成解密种子失败：%v", err)
      		return
      	}
      	sm2Decrypter = sm2.New(semDecrypter) // 生成解密种子
      	if sm4Decrypter, err = sm4.New(sm4.KeyBytes(key), sm4.IVBytes(iv)); err != nil {
      		log.Fatalf("生成解密器(SM4)失败：%v", err)
      	}
      	if err = sm4Decrypter.DecryptCBCLargeFile(encryptedFile, decryptedFile, sm2Decrypter); err != nil {
      		log.Fatalf("大文件解密失败：%v", err)
      	}
      	log.Printf("大文件解密成功：%s", decryptedFile)
      
      	// 4) 流式对比，避免将大文件全部读入内存
      	func() {
      		var (
      			f1, f2 *os.File
      			b1     = make([]byte, 1024*1024)
      			b2     = make([]byte, 1024*1024)
      		)
      
      		if f1, err = os.Open(plainFile); err != nil {
      			return
      		}
      		defer func() { _ = f1.Close() }()
      
      		if f2, err = os.Open(decryptedFile); err != nil {
      			return
      		}
      		defer func() { _ = f2.Close() }()
      
      		for {
      			n1, e1 := io.ReadFull(f1, b1)
      			n2, e2 := io.ReadFull(f2, b2)
      			if n1 != n2 {
      				err = errors.New("文件长度不一致")
      				return
      			}
      			for i := range n1 {
      				if b1[i] != b2[i] {
      					err = errors.New("文件内容不一致")
      					return
      				}
      			}
      			if e1 == io.EOF || e1 == io.ErrUnexpectedEOF {
      				break
      			}
      			if e1 != nil || e2 != nil {
      				err = errors.New("文件对比失败")
      				return
      			}
      		}
      	}()
      
      	if err != nil {
      		log.Fatalf("大文件对比失败：%v", err)
      	}
      
      	log.Print("✓ 大文件内容一致，加解密正确")
      	cleanupTestFiles(plainFile, encryptedFile, decryptedFile)
      }
      
      func cleanupTestFiles(files ...string) {
      	for _, file := range files {
      		if rmErr := os.Remove(file); rmErr != nil && !os.IsNotExist(rmErr) {
      			log.Printf("清理测试文件失败 %s: %v", file, rmErr)
      		}
      	}
      }
      
      func main() { TestFileEncrypt(); TestLargeFileEncrypt() }
      ```
   
   2. *RSA+AES128-CBC*
   
      ```go
      package main
      
      import (
      	"crypto/rand"
      	"errors"
      	"io"
      	"log"
      	"os"
      
      	"github.com/aid297/aid/v2/secret"
      	"github.com/aid297/aid/v2/secret/asymmetric/rsa"
      	"github.com/aid297/aid/v2/secret/symmetric/aes"
      )
      
      // 1. 文件加密/解密演示（RSA + AES）
      func TestFileEncrypt() {
      	var (
      		err                        error
      		semEncrypter, semDecrypter secret.Semener
      		rsaEncrypter, rsaDecrypter secret.Asymmetricer
      		aesEncrypter, aesDecrypter secret.Symmetricer
      		aesKey, aesIV, priKeyBytes []byte
      		plainFile                  = "/tmp/rsa_test_plain.txt"
      		encryptedFile              = "/tmp/rsa_test_encrypted.bin"
      		decryptedFile              = "/tmp/rsa_test_decrypted.txt"
      	)
      
      	if semEncrypter, err = rsa.NewSem(); err != nil {
      		log.Fatalf("生成加密种子(RSA)失败：%v", err)
      	}
      	rsaEncrypter = rsa.New(semEncrypter)
      
      	if err = os.WriteFile(plainFile, []byte("这是需要加密的文件内容（RSA + AES），可以是任意大小的数据。"), 0644); err != nil {
      		log.Fatalf("写入测试文件失败：%v", err)
      	}
      
      	if aesEncrypter, err = aes.New(aes.RandKeyWithBits(aes.AESKey192, &aesKey), aes.RandIV(&aesIV)); err != nil {
      		log.Fatalf("创建加密器(AES)失败：%v", err)
      	}
      	if err = aesEncrypter.EncryptCBCFile(plainFile, encryptedFile, rsaEncrypter); err != nil {
      		log.Fatalf("加密失败：%v", err)
      	}
      	log.Printf("加密成功：%s", encryptedFile)
      
      	if priKeyBytes, err = semEncrypter.GetPriKeyBytes(); err != nil {
      		log.Fatalf("获取加密种子私钥失败：%v", err)
      	}
      
      	if semDecrypter, err = rsa.NewSem(rsa.PriKeyBytes(priKeyBytes)); err != nil {
      		log.Fatalf("生成解密种子(RSA)失败：%v", err)
      	}
      	rsaDecrypter = rsa.New(semDecrypter)
      
      	if aesDecrypter, err = aes.New(aes.KeyBytes(aesKey), aes.IVBytes(aesIV)); err != nil {
      		log.Fatalf("生成解密器(AES)失败：%v", err)
      	}
      	if err = aesDecrypter.DecryptCBCFile(encryptedFile, decryptedFile, rsaDecrypter); err != nil {
      		log.Fatalf("解密文件失败：%v", err)
      	}
      	log.Printf("解密成功：%s", decryptedFile)
      
      	original, _ := os.ReadFile(plainFile)
      	decrypted, _ := os.ReadFile(decryptedFile)
      	if string(original) == string(decrypted) {
      		log.Print("✓ 文件内容一致，加解密正确")
      	} else {
      		log.Fatal("✗ 文件内容不一致")
      	}
      
      	cleanupTestFiles(plainFile, encryptedFile, decryptedFile)
      }
      
      // 2. 大文件加密/解密演示（流式，RSA + AES）
      func TestLargeFileEncrypt() {
      	var (
      		err                        error
      		semEncrypter, semDecrypter secret.Semener
      		rsaEncrypter, rsaDecrypter secret.Asymmetricer
      		aesEncrypter, aesDecrypter secret.Symmetricer
      		aesKey, aesIV, priKeyBytes []byte
      		plainFile                        = "/tmp/rsa_test_large_plain.bin"
      		encryptedFile                    = "/tmp/rsa_test_large_encrypted.bin"
      		decryptedFile                    = "/tmp/rsa_test_large_decrypted.bin"
      		fileSize                   int64 = 64 * 1024 * 1024 // 64MB 演示
      	)
      
      	if semEncrypter, err = rsa.NewSem(); err != nil {
      		log.Fatalf("生成加密种子(RSA)失败：%v", err)
      	}
      	rsaEncrypter = rsa.New(semEncrypter)
      
      	func() {
      		var (
      			f   *os.File
      			buf = make([]byte, 1024*1024)
      			wr  int64
      		)
      		if f, err = os.Create(plainFile); err != nil {
      			return
      		}
      		defer f.Close()
      
      		for wr < fileSize {
      			need := min(fileSize-wr, int64(len(buf)))
      			if _, err = rand.Read(buf[:need]); err != nil {
      				return
      			}
      			if _, err = f.Write(buf[:need]); err != nil {
      				return
      			}
      			wr += need
      		}
      	}()
      	if err != nil {
      		log.Fatalf("生成大文件失败：%v", err)
      	}
      
      	if aesEncrypter, err = aes.New(aes.RandKeyWithBits(aes.AESKey256), aes.RandIV()); err != nil {
      		log.Fatalf("生成 AES 失败：%v", err)
      	}
      	if err = aesEncrypter.EncryptCBCLargeFile(plainFile, encryptedFile, rsaEncrypter); err != nil {
      		log.Fatalf("大文件加密失败：%v", err)
      	}
      	log.Printf("大文件加密成功：%s", encryptedFile)
      
      	if priKeyBytes, err = semEncrypter.GetPriKeyBytes(); err != nil {
      		log.Fatalf("获取加密种子私钥失败：%v", err)
      	}
      
      	if semDecrypter, err = rsa.NewSem(rsa.PriKeyBytes(priKeyBytes)); err != nil {
      		log.Fatalf("生成解密种子(RSA)失败：%v", err)
      	}
      	rsaDecrypter = rsa.New(semDecrypter)
      
      	if aesDecrypter, err = aes.New(aes.KeyBytes(aesKey), aes.IVBytes(aesIV)); err != nil {
      		log.Fatalf("生成解密器(AES)失败：%v", err)
      	}
      
      	if err = aesDecrypter.DecryptCBCLargeFile(encryptedFile, decryptedFile, rsaDecrypter); err != nil {
      		log.Fatalf("大文件解密失败：%v", err)
      	}
      	log.Printf("大文件解密成功：%s", decryptedFile)
      
      	func() {
      		var (
      			f1, f2 *os.File
      			b1     = make([]byte, 1024*1024)
      			b2     = make([]byte, 1024*1024)
      		)
      
      		if f1, err = os.Open(plainFile); err != nil {
      			return
      		}
      		defer f1.Close()
      
      		if f2, err = os.Open(decryptedFile); err != nil {
      			return
      		}
      		defer f2.Close()
      
      		for {
      			n1, e1 := io.ReadFull(f1, b1)
      			n2, e2 := io.ReadFull(f2, b2)
      			if n1 != n2 {
      				err = errors.New("文件长度不一致")
      				return
      			}
      			for i := range n1 {
      				if b1[i] != b2[i] {
      					err = errors.New("文件内容不一致")
      					return
      				}
      			}
      			if e1 == io.EOF || e1 == io.ErrUnexpectedEOF {
      				break
      			}
      			if e1 != nil || e2 != nil {
      				err = errors.New("文件对比失败")
      				return
      			}
      		}
      	}()
      
      	if err != nil {
      		log.Fatalf("大文件对比失败：%v", err)
      	}
      
      	log.Print("✓ 大文件内容一致，加解密正确")
      	cleanupTestFiles(plainFile, encryptedFile, decryptedFile)
      }
      
      func cleanupTestFiles(files ...string) {
      	for _, file := range files {
      		if rmErr := os.Remove(file); rmErr != nil && !os.IsNotExist(rmErr) {
      			log.Printf("清理测试文件失败 %s: %v", file, rmErr)
      		}
      	}
      }
      
      func main() { TestFileEncrypt(); TestLargeFileEncrypt() }
      ```