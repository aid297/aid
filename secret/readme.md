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
      		err                        error
      		semEncrypter, semDecrypter secret.Semener
      		sm2Encrypter, sm2Decrypter secret.Asymmetricer
      		sm2SemAPriKeyBytes         []byte
      		plainText                  = []byte("hello, SM2 非对称加密测试!")
      		cipherBase64               string
      		decrypted                  []byte
      	)
      
      	if semEncrypter, err = sm2.NewSem(); err != nil {
      		t.Fatalf("生成加密种子失败：%v", err)
      	}
      
      	sm2Encrypter = sm2.New(semEncrypter)
      
      	if cipherBase64, err = sm2Encrypter.Encrypt(plainText); err != nil {
      		t.Fatalf("加密错误：%v", err)
      	}
      	t.Logf("加密结果：%s\n", cipherBase64)
      
      	if sm2SemAPriKeyBytes, err = semEncrypter.GetPriKeyBytes(); err != nil {
      		t.Fatalf("获取私钥失败：%v", err)
      	}
      
      	if semDecrypter, err = sm2.NewSem(sm2.PriKeyBytes(sm2SemAPriKeyBytes)); err != nil {
      		t.Fatalf("生成解密种子失败：%v", err)
      	}
      	a, _ := semDecrypter.GetPriKeyBase64()
      	b, _ := semDecrypter.GetPubKeyBase64()
      	t.Logf("解密种子私钥：%s, 公钥：%s", a, b)
      
      	sm2Decrypter = sm2.New(semDecrypter)
      
      	if decrypted, err = sm2Decrypter.Decrypt(cipherBase64); err != nil {
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
      		err                error
      		semSign, semVerify secret.Semener
      		sm2Sign, sm2Verify secret.Asymmetricer
      		data               = []byte("hello, SM2 数字签名测试!")
      		sigHex             string
      		ok                 bool
      	)
      
      	if semSign, err = sm2.NewSem(); err != nil {
      		t.Fatalf("生成签名种子失败：%v", err)
      	}
      	sm2Sign = sm2.New(semSign)
      
      	if sigHex, err = sm2Sign.Sign(data); err != nil {
      		t.Fatalf("签名失败：%v", err)
      	}
      	t.Logf("签名内容(hex): %s", sigHex)
      
      	if semVerify, err = sm2.NewSem(sm2.PriKey(semSign.GetPriKey())); err != nil {
      		t.Fatalf("生成验证种子失败：%v", err)
      	}
      	sm2Verify = sm2.New(semVerify)
      
      	if ok, err = sm2Verify.Verify(data, sigHex); err != nil {
      		t.Fatalf("验证失败：%v", err)
      	}
      	if !ok {
      		t.Fatal("签名验证失败")
      	}
      
      	t.Logf("验证成功")
      }
      
      func TestVerifyWithWrongData(t *testing.T) {
      	var (
      		err                error
      		semSign, semVerify secret.Semener
      		sm2Sign, sm2Verify secret.Asymmetricer
      		data               = []byte("original data")
      		sigHex             string
      		ok                 bool
      	)
      
      	if semSign, err = sm2.NewSem(); err != nil {
      		t.Fatalf("生成种子失败：%v", err)
      	}
      	sm2Sign = sm2.New(semSign)
      
      	if sigHex, err = sm2Sign.Sign(data); err != nil {
      		t.Fatalf("签名失败：%v", err)
      	}
      
      	if semVerify, err = sm2.NewSem(sm2.PriKey(semSign.GetPriKey())); err != nil {
      		t.Fatalf("生成解密种子失败：%v", err)
      	}
      	sm2Verify = sm2.New(semVerify)
      
      	if ok, err = sm2Verify.Verify([]byte("tampered data"), sigHex); err != nil {
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
      		err          error
      		semA, semB   secret.Semener
      		rsaA, rsaB   secret.Asymmetricer
      		plainText    = []byte("hello, RSA 非对称加密测试!")
      		cipherBase64 string
      		decrypted    []byte
      	)
      
      	if semA, err = rsa.NewSem(); err != nil {
      		t.Fatalf("生成加密种子失败：%v", err)
      	}
      	rsaA = rsa.New(semA)
      
      	if cipherBase64, err = rsaA.Encrypt(plainText); err != nil {
      		t.Fatalf("加密错误：%v", err)
      	}
      	t.Logf("加密结果：%s\n", cipherBase64)
      
      	if semB, err = rsa.NewSem(rsa.PriKey(semA.GetPriKey())); err != nil {
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
      		err                error
      		semSign, semVerify secret.Semener
      		rsaSign, rsaVerify secret.Asymmetricer
      		data               = []byte("hello, RSA 数字签名测试!")
      		sigHex             string
      		ok                 bool
      	)
      
      	if semSign, err = rsa.NewSem(); err != nil {
      		t.Fatalf("生成签名种子失败：%v", err)
      	}
      	rsaSign = rsa.New(semSign)
      
      	if sigHex, err = rsaSign.Sign(data); err != nil {
      		t.Fatalf("签名失败：%v", err)
      	}
      	t.Logf("签名内容(hex): %s", sigHex)
      
      	if semVerify, err = rsa.NewSem(rsa.PriKey(semSign.GetPriKey())); err != nil {
      		t.Fatalf("生成验证种子失败：%v", err)
      	}
      	rsaVerify = rsa.New(semVerify)
      
      	if ok, err = rsaVerify.Verify(data, sigHex); err != nil {
      		t.Fatalf("验证失败：%v", err)
      	}
      	if !ok {
      		t.Fatal("签名验证失败")
      	}
      
      	t.Logf("验证成功")
      }
      
      func TestVerifyWithWrongData(t *testing.T) {
      	var (
      		err                error
      		semSign, semVerify secret.Semener
      		rsaSign, rsaVerify secret.Asymmetricer
      		data               = []byte("original data")
      		sigHex             string
      		ok                 bool
      	)
      	if semSign, err = rsa.NewSem(); err != nil {
      		t.Fatalf("生成签名种子失败：%v", err)
      	}
      
      	if err = rsa.MustGeneratePriKey(semSign); err != nil {
      		t.Fatalf("%v", err)
      	}
      
      	rsaSign = rsa.New(semSign)
      	if sigHex, err = rsaSign.Sign(data); err != nil {
      		t.Fatalf("签名失败：%v", err)
      	}
      
      	if semVerify, err = rsa.NewSem(rsa.PriKey(semSign.GetPriKey())); err != nil {
      		t.Fatalf("生成验证种子失败：%v", err)
      	}
      	rsaVerify = rsa.New(semVerify)
      
      	if ok, err = rsaVerify.Verify([]byte("tampered data"), sigHex); err != nil {
      		t.Fatalf("验证失败：%v", err)
      	}
      	if ok {
      		t.Fatal("篡改数据的验证应失败")
      	}
      	t.Logf("篡改数据正确被拒绝")
      }
      ```
   
3. 对称加密（*SM4*）

   *SM4* 支持 *ECB* 和 *CBC* 两种模式，通过 `AlgorithmCBC()`/`AlgorithmECB()` 选择。默认 *CBC* 模式。

   1. 基本用法（默认 *CBC* 模式）
      ```go
      func TestCBC(t *testing.T) {
      	sm4Helper, err := sm4.New(sm4.KeyBytes(testKey), sm4.IVBytes(testIV))
      	if err != nil {
      		t.Fatalf("创建 SM4 对象失败：%v", err)
      	}
      	cipherText, err := sm4Helper.Encrypt(testPlain)
      	if err != nil {
      		t.Fatalf("Encrypt failed: %v", err)
      	}
      	plain, err := sm4Helper.Decrypt(cipherText)
      	if err != nil {
      		t.Fatalf("Decrypt failed: %v", err)
      	}
      	if !bytes.Equal(plain, testPlain) {
      		t.Fatalf("SM4 mismatch: got %s, want %s", plain, testPlain)
      	}
      	t.Logf("SM4 OK: %s", plain)
      }
      ```

   2. *ECB* 模式：由于*ECB*模式不安全，所以不推荐
      
      ```go
      func TestECB(t *testing.T) {
      	sm4Helper, err := sm4.New(sm4.KeyBytes(testKey), sm4.AlgorithmECB())
      	if err != nil {
      		t.Fatalf("创建 ECB 对象失败：%v", err)
      	}
      	cipherText, err := sm4Helper.Encrypt(testPlain)
      	if err != nil {
      		t.Fatalf("Encrypt failed: %v", err)
      	}
      	plain, err := sm4Helper.Decrypt(cipherText)
      	if err != nil {
      		t.Fatalf("Decrypt failed: %v", err)
      	}
      	if !bytes.Equal(plain, testPlain) {
      		t.Fatalf("ECB mismatch: got %s, want %s", plain, testPlain)
      	}
      	t.Logf("ECB OK: %s", plain)
      }
      ```
      
   3. Base64 编解码
      ```go
      func TestECBBase64(t *testing.T) {
      	sm4Helper, err := sm4.New(sm4.KeyBytes(testKey))
      	if err != nil {
      		t.Fatalf("创建 SM4 对象失败：%v", err)
      	}
      	b64, err := sm4Helper.EncryptBase64(testPlain)
      	if err != nil {
      		t.Fatalf("加密失败: %v", err)
      	}
      	t.Logf("Base64: %s", b64)
      	plain, err := sm4Helper.DecryptBase64(b64)
      	if err != nil {
      		t.Fatalf("解密失败: %v", err)
      	}
      	if !bytes.Equal(plain, testPlain) {
      		t.Fatalf("匹配失败: 结果 %s，期望 %s", plain, testPlain)
      	}
      	t.Logf("Base64 OK: %s", plain)
      }
      ```

5. 对称加密（*AES*）

   AES 支持 ECB 和 CBC 两种模式，通过 `AlgorithmCBC()`/`AlgorithmECB()` 选择。默认 CBC 模式。

   1. 基本用法（默认 CBC 模式）
      ```go
      func TestCBC(t *testing.T) {
      	var (
      		testKey   = []byte("1234567890abcdef")
      		testIV    = []byte("abcdef1234567890")
      		testPlain = []byte("hello, AES encrypt test!")
      	)
      
      	aesHelper, err := aes.New(aes.KeyBytes(testKey), aes.IVBytes(testIV))
      	if err != nil {
      		t.Fatalf("创建 AES 对象失败：%v", err)
      	}
      
      cipherText, err := aesHelper.Encrypt(testPlain)
      	if err != nil {
      		t.Fatalf("加密失败：%v", err)
      	}
      
      	plain, err := aesHelper.Decrypt(cipherText)
      	if err != nil {
      		t.Fatalf("解密失败：%v", err)
      	}
      
      	if !bytes.Equal(plain, testPlain) {
      		t.Fatalf("比对失败：结果 %s，期望 %s", plain, testPlain)
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
      		t.Fatalf("创建 AES 对象失败：%v", err)
      	}
      
      b64, err := aesHelper.EncryptBase64(testPlain)
      	if err != nil {
      		t.Fatalf("加密失败：%v", err)
      	}
      
      	plain, err := aesHelper.DecryptBase64(b64)
      	if err != nil {
      		t.Fatalf("解密失败：%v", err)
      	}
      
      	if !bytes.Equal(plain, testPlain) {
      		t.Fatalf("比对失败：结果 %s，期望 %s", plain, testPlain)
      	}
      }
      ```

   2. ECB 模式
      ```go
      func TestECB(t *testing.T) {
      	var (
      		testKey   = []byte("1234567890abcdef")
      		testPlain = []byte("hello, AES encrypt test!")
      	)
      
      	aesHelper, err := aes.New(aes.KeyBytes(testKey), aes.AlgorithmECB())
      	if err != nil {
      		t.Fatalf("创建 ECB 对象失败：%v", err)
      	}
      
      cipherText, err := aesHelper.Encrypt(testPlain)
      	if err != nil {
      		t.Fatalf("加密失败：%v", err)
      	}
      
      	plain, err := aesHelper.Decrypt(cipherText)
      	if err != nil {
      		t.Fatalf("解密失败：%v", err)
      	}
      
      	if !bytes.Equal(plain, testPlain) {
      		t.Fatalf("比对失败：结果 %s，期望 %s", plain, testPlain)
      	}
      }
      
      func TestECBBase64(t *testing.T) {
      	var (
      		testKey   = []byte("1234567890abcdef")
      		testPlain = []byte("hello, AES encrypt test!")
      	)
      
      	aesHelper, err := aes.New(aes.KeyBytes(testKey), aes.AlgorithmECB())
      	if err != nil {
      		t.Fatalf("创建 ECB 对象失败：%v", err)
      	}
      
      b64, err := aesHelper.EncryptBase64(testPlain)
      	if err != nil {
      		t.Fatalf("加密失败：%v", err)
      	}
      
      	plain, err := aesHelper.DecryptBase64(b64)
      	if err != nil {
      		t.Fatalf("解密失败：%v", err)
      	}
      
      	if !bytes.Equal(plain, testPlain) {
      		t.Fatalf("比对失败：结果 %s，期望 %s", plain, testPlain)
      	}
      }
      ```

   3. AES192 和 AES256 支持
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
      		cipherText, err := aesHelper.Encrypt(testPlain)
      			if err != nil {
      				t.Fatalf("加密失败：%v", err)
      			}
      			plain, err := aesHelper.Decrypt(cipherText)
      			if err != nil {
      				t.Fatalf("解密失败：%v", err)
      			}
      			if !bytes.Equal(plain, testPlain) {
      				t.Fatalf("比对失败：结果 %s，期望 %s", plain, testPlain)
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
      
      if _, err = aesHelper.Encrypt(testPlain); err != nil {
      		t.Fatalf("加密失败：%v", err)
      	}
      }
      ```

7. 组合用法

   1. *SM2+SM4*加密文件和大文件加密

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
      
      	// 加密（默认 CBC 模式）
      	if sm4Encrypter, err = sm4.New(sm4.RandKey(&key), sm4.RandIV(&iv)); err != nil {
      		log.Fatalf("生成加密器(SM4)失败：%v", err)
      	}
      	if err = sm4Encrypter.EncryptFile(plainFile, encryptedFile, sm2Encrypter); err != nil {
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
      	if err = sm4Decrypter.DecryptFile(encryptedFile, decryptedFile, sm2Decrypter); err != nil {
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
      		plainFile                  = "/tmp/sm2_test_large_plain.bin"
      		encryptedFile              = "/tmp/sm2_test_large_encrypted.bin"
      		decryptedFile              = "/tmp/sm2_test_large_decrypted.bin"
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
      
      	// 2) 流式加密（默认 CBC 模式）
      	if sm4Encrypter, err = sm4.New(sm4.RandKey(&key), sm4.RandIV(&iv)); err != nil {
      		log.Fatalf("生成加密器(SM4)失败：%v", err)
      	}
      	if err = sm4Encrypter.EncryptLargeFile(plainFile, encryptedFile, sm2Encrypter); err != nil {
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
      	if err = sm4Decrypter.DecryptLargeFile(encryptedFile, decryptedFile, sm2Decrypter); err != nil {
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
   
   2. *RSA+AES*
   
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
      	if err = aesEncrypter.EncryptFile(plainFile, encryptedFile, rsaEncrypter); err != nil {
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
      	if err = aesDecrypter.DecryptFile(encryptedFile, decryptedFile, rsaDecrypter); err != nil {
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
      		plainFile                  = "/tmp/rsa_test_large_plain.bin"
      		encryptedFile              = "/tmp/rsa_test_large_encrypted.bin"
      		decryptedFile              = "/tmp/rsa_test_large_decrypted.bin"
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
      	if err = aesEncrypter.EncryptLargeFile(plainFile, encryptedFile, rsaEncrypter); err != nil {
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
      
      	if err = aesDecrypter.DecryptLargeFile(encryptedFile, decryptedFile, rsaDecrypter); err != nil {
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

8. 统一接口说明

   对称加密接口支持通过选项函数选择 ECB 或 CBC 模式：

   - `sm4.AlgorithmECB()` / `aes.AlgorithmECB()` - 电子密码本模式（各区块独立加密，不推荐用于加密敏感数据）
   - `sm4.AlgorithmCBC()` / `aes.AlgorithmCBC()` - 密码块链接模式（需配合 IV 使用，安全性更高）
   - **默认模式为 CBC**

   **统一方法列表：**
   | 方法 | 说明 |
   |------|------|
   | `Encrypt(in, out io.Reader/Writer)` | 流式加解密 |
   | `EncryptBase64/DecryptBase64` | Base64 编解码 |
   | `EncryptFile/DecryptFile` | 文件加解密 |
   | `EncryptLargeFile/DecryptLargeFile` | 大文件流式加解密 |