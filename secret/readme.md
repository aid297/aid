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
      
      	if sem, err = NewSem(); err != nil {
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
      		sm2A, sm2B         secret.Asymmetricor
      		sm2SemAPriKeyBytes []byte
      		plainText          = []byte("hello, SM2 非对称加密测试!")
      		cipherBase64       string
      		decrypted          []byte
      	)
      
      	if semA, err = NewSem(); err != nil {
      		t.Fatalf("生成加密种子失败：%v", err)
      	}
      
      	sm2A = New(semA)
      
      	if cipherBase64, err = sm2A.Encrypt(plainText); err != nil {
      		t.Fatalf("加密错误：%v", err)
      	}
      	t.Logf("加密结果：%s\n", cipherBase64)
      
      	if sm2SemAPriKeyBytes, err = semA.GetPriKeyBytes(); err != nil {
      		t.Fatalf("获取私钥失败：%v", err)
      	}
      
      	if semB, err = NewSem(PriKeyBytes(sm2SemAPriKeyBytes)); err != nil {
      		t.Fatalf("生成解密种子失败：%v", err)
      	}
      	a, _ := semB.GetPriKeyBase64()
      	b, _ := semB.GetPubKeyBase64()
      	t.Logf("解密种子私钥：%s, 公钥：%s", a, b)
      
      	sm2B = New(semB)
      
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
      		err    error
      		semA   secret.Semener
      		sm2    secret.Asymmetricor
      		data   = []byte("hello, SM2 数字签名测试!")
      		sigHex string
      		ok     bool
      	)
      
      	if semA, err = NewSem(); err != nil {
      		t.Fatalf("生成种子失败：%v", err)
      	}
      
      	sm2 = New(semA)
      
      	if sigHex, err = sm2.Sign(data); err != nil {
      		t.Fatalf("签名失败：%v", err)
      	}
      	t.Logf("签名内容(hex): %s", sigHex)
      
      	if ok, err = sm2.Verify(data, sigHex); err != nil {
      		t.Fatalf("验证失败：%v", err)
      	}
      	if !ok {
      		t.Fatal("签名验证失败")
      	}
      
      	t.Logf("验证成功")
      }
      ```

2. 非对称加密（*RSA*）

   1. 基本用法
      ```go
      func TestEncryptDecrypt(t *testing.T) {
      	pub, pri, err := GenerateKeyPairBase64(2048)
      	if err != nil {
      		t.Fatalf("generate key failed: %v", err)
      	}
      
      	helper, err := NewByBase64(pub, pri)
      	if err != nil {
      		t.Fatalf("new helper failed: %v", err)
      	}
      
      	plain := []byte("hello, RSA interface test!")
      	cipherBase64, err := helper.Encrypt(plain)
      	if err != nil {
      		t.Fatalf("encrypt failed: %v", err)
      	}
      
      	decrypted, err := helper.Decrypt(cipherBase64)
      	if err != nil {
      		t.Fatalf("decrypt failed: %v", err)
      	}
      	if !bytes.Equal(plain, decrypted) {
      		t.Fatalf("decrypt mismatch: got %s, want %s", decrypted, plain)
      	}
      }
      ```

   2. 签名&验证
      ```go
      func TestSignVerify(t *testing.T) {
      	pub, pri, err := GenerateKeyPairBase64(2048)
      	if err != nil {
      		t.Fatalf("generate key failed: %v", err)
      	}
      
      	helper, err := NewByBase64(pub, pri)
      	if err != nil {
      		t.Fatalf("new helper failed: %v", err)
      	}
      
      	data := []byte("hello, rsa sign")
      	sig, err := helper.Sign(data)
      	if err != nil {
      		t.Fatalf("sign failed: %v", err)
      	}
      
      	ok, err := helper.Verify(data, sig)
      	if err != nil {
      		t.Fatalf("verify failed: %v", err)
      	}
      	if !ok {
      		t.Fatalf("verify should be true")
      	}
      
      	ok, err = helper.Verify([]byte("tampered"), sig)
      	if err != nil {
      		t.Fatalf("verify tampered failed: %v", err)
      	}
      	if ok {
      		t.Fatalf("verify should be false for tampered data")
      	}
      }
      ```

      

3. 对称加密（*SM4-ECB*）

   1. 基本用法
      ```go
      func TestECB(t *testing.T) {
      	sm4Helper, err := New(KeyBytes(testKey))
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
      
      
      func TestECB(t *testing.T) {
      	sm4Helper, err := New(KeyBytes(testKey))
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
      ```

4. 对称加密（*SM4-CBC*）

   1. 基本用法
      ```go
      func TestCBC(t *testing.T) {
      	sm4Helper, err := New(KeyBytes(testKey), IVBytes(testIV))
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
      	sm4Helper, err := New(KeyBytes(testKey), IVBytes(testIV))
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
      	sm4Helper, err := New(KeyString("shortkey"))
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
      
      	aesHelper, err := New(KeyBytes(testKey))
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
      
      	aesHelper, err := New(KeyBytes(testKey))
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
      
      	aesHelper, err := New(KeyBytes(testKey), IVBytes(testIV))
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
      
      	aesHelper, err := New(KeyBytes(testKey), IVBytes(testIV))
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
      
      	aesHelper, err := New(KeyString("shortkey"))
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
      			aesHelper, err := New(KeyBytes(tt.key), IVBytes(testIV))
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
      
      	aesHelper, err := New(RandKeyWithBits(KeyBits256, &key), RandIV())
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
      
      	aesHelper, err := New(KeySize(KeyBits192), RandKey(&key), RandIV())
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
      	_, err := New(RandKeyWithBits(111), RandIV())
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
      
      // 1. 文件加密/解密演示
      func TestCBCEncryptDecryptFile() {
      	// 生成密钥对
      	var (
      		err           error
      		sem           secret.Semener
      		sm2Helper     secret.Asymmetricer
      		plainFile     = "/tmp/sm2_test_plain.txt"
      		encryptedFile = "/tmp/sm2_test_encrypted.bin"
      		decryptedFile = "/tmp/sm2_test_decrypted.txt"
      		sm4Helper     secret.Symmetricer
      		sm4Key        []byte
      		sm4IV         []byte
      	)
      
      	if sem, err = sm2.NewSem(); err != nil {
      		log.Fatalf("生成种子失败：%v", err)
      	}
      
      	sm2Helper = sm2.New(sem)
      
      	// 准备测试文件
      	if err = os.WriteFile(plainFile, []byte("这是需要加密的文件内容，可以是任意大小的数据。"), 0644); err != nil {
      		log.Fatalf("写入测试文件失败：%v", err)
      	}
      
      	// 加密
      	if sm4Helper, err = sm4.New(sm4.RandKey(&sm4Key), sm4.RandIV(&sm4IV)); err != nil {
      		log.Fatalf("生成 SM4 失败：%v", err)
      	}
      	if err = sm4Helper.EncryptCBCFile(plainFile, encryptedFile, sm2Helper); err != nil {
      		log.Fatalf("加密失败：%v", err)
      	}
      	log.Printf("加密成功：%s", encryptedFile)
      
      	// 解密
      	if err = sm4Helper.DecryptCBCFile(encryptedFile, decryptedFile, sm2Helper); err != nil {
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
      
      // 2. 大文件加密/解密演示（流式）
      func TestCBCLargeFileEncryptDecrypt() {
      	var (
      		err           error
      		sem           secret.Semener
      		sm2Helper     secret.Asymmetricer
      		plainFile           = "/tmp/sm2_test_large_plain.bin"
      		encryptedFile       = "/tmp/sm2_test_large_encrypted.bin"
      		decryptedFile       = "/tmp/sm2_test_large_decrypted.bin"
      		fileSize      int64 = 64 * 1024 * 1024 // 64MB 演示
      		sm4Helper     secret.Symmetricer
      	)
      
      	if sem, err = sm2.NewSem(); err != nil {
      		log.Fatalf("生成种子失败：%v", err)
      	}
      
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
      
      	if sm4Helper, err = sm4.New(sm4.RandKey(), sm4.RandIV()); err != nil {
      		log.Fatalf("生成 SM4 失败：%v", err)
      	}
      
      	// 2) 流式加密
      	sm2Helper = sm2.New(sem)
      	if err = sm4Helper.EncryptCBCLargeFile(plainFile, encryptedFile, sm2Helper); err != nil {
      		log.Fatalf("大文件加密失败：%v", err)
      	}
      	log.Printf("大文件加密成功：%s", encryptedFile)
      
      	// 3) 流式解密
      	if err = sm4Helper.DecryptCBCLargeFile(encryptedFile, decryptedFile, sm2Helper); err != nil {
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
      
      func main() { TestCBCLargeFileEncryptDecrypt() }
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
      func TestCBCEncryptDecryptFile() {
      	var (
      		err           error
      		rsaHelper     secret.Asymmetricer
      		plainFile     = "/tmp/rsa_test_plain.txt"
      		encryptedFile = "/tmp/rsa_test_encrypted.bin"
      		decryptedFile = "/tmp/rsa_test_decrypted.txt"
      		aesHelper     secret.Symmetricer
      		aesKey        []byte
      		aesIV         []byte
      	)
      
      	pubBase64, priBase64, err := rsa.GenerateKeyPairBase64(2048)
      	if err != nil {
      		log.Fatalf("生成 RSA 密钥对失败：%v", err)
      	}
      
      	rsaHelper, err = rsa.NewByBase64(pubBase64, priBase64)
      	if err != nil {
      		log.Fatalf("创建 RSA 实例失败：%v", err)
      	}
      
      	if err = os.WriteFile(plainFile, []byte("这是需要加密的文件内容（RSA + AES），可以是任意大小的数据。"), 0644); err != nil {
      		log.Fatalf("写入测试文件失败：%v", err)
      	}
      
      	if aesHelper, err = aes.New(aes.RandKey(&aesKey), aes.RandIV(&aesIV)); err != nil {
      		log.Fatalf("生成 AES 失败：%v", err)
      	}
      	if err = aesHelper.EncryptCBCFile(plainFile, encryptedFile, rsaHelper); err != nil {
      		log.Fatalf("加密失败：%v", err)
      	}
      	log.Printf("加密成功：%s", encryptedFile)
      
      	if err = aesHelper.DecryptCBCFile(encryptedFile, decryptedFile, rsaHelper); err != nil {
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
      func TestCBCLargeFileEncryptDecrypt() {
      	var (
      		err           error
      		rsaHelper     secret.Asymmetricer
      		plainFile           = "/tmp/rsa_test_large_plain.bin"
      		encryptedFile       = "/tmp/rsa_test_large_encrypted.bin"
      		decryptedFile       = "/tmp/rsa_test_large_decrypted.bin"
      		fileSize      int64 = 64 * 1024 * 1024 // 64MB 演示
      		aesHelper     secret.Symmetricer
      	)
      
      	pubBase64, priBase64, err := rsa.GenerateKeyPairBase64(2048)
      	if err != nil {
      		log.Fatalf("生成 RSA 密钥对失败：%v", err)
      	}
      
      	rsaHelper, err = rsa.NewByBase64(pubBase64, priBase64)
      	if err != nil {
      		log.Fatalf("创建 RSA 实例失败：%v", err)
      	}
      
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
      
      	if aesHelper, err = aes.New(aes.RandKey(), aes.RandIV()); err != nil {
      		log.Fatalf("生成 AES 失败：%v", err)
      	}
      
      	if err = aesHelper.EncryptCBCLargeFile(plainFile, encryptedFile, rsaHelper); err != nil {
      		log.Fatalf("大文件加密失败：%v", err)
      	}
      	log.Printf("大文件加密成功：%s", encryptedFile)
      
      	if err = aesHelper.DecryptCBCLargeFile(encryptedFile, decryptedFile, rsaHelper); err != nil {
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
      
      func main() { TestCBCLargeFileEncryptDecrypt() }
      ```

   3. 