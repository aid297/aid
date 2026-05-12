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

   3. 验证签名
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

   4. 验证篡改数据
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

2. 非对称加密（*SM4-ECB*）

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

3. 非对称加密（*SM4-CBC*）

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

   2. 加解密文件（组合*SM2+SM4-CBC*）
      ```go
      // 文件加密/解密演示
      func TestCBCEncryptDecryptFile() {
      	// 生成密钥对
      	var (
      		err           error
      		sem           secret.Semener
      		plainFile     = "/tmp/sm2_test_plain.txt"
      		encryptedFile = "/tmp/sm2_test_encrypted.bin"
      		decryptedFile = "/tmp/sm2_test_decrypted.txt"
      		sm4Helper     secret.Symmetricor
      		sm4Key        []byte
      		sm4IV         []byte
      	)
      
      	if sem, err = sm2.NewSem(); err != nil {
      		log.Fatalf("生成种子失败：%v", err)
      	}
      
      	// 准备测试文件
      	if err = os.WriteFile(plainFile, []byte("这是需要加密的文件内容，可以是任意大小的数据。"), 0644); err != nil {
      		log.Fatalf("写入测试文件失败：%v", err)
      	}
      
      	// 加密
      	if sm4Helper, err = sm4.New(sm4.RandKey(&sm4Key), sm4.RandIV(&sm4IV)); err != nil {
      		log.Fatalf("生成 SM4 失败：%v", err)
      	}
      	if err = sm4Helper.EncryptCBCFile(plainFile, encryptedFile, sem); err != nil {
      		log.Fatalf("加密失败：%v", err)
      	}
      	log.Printf("加密成功：%s", encryptedFile)
      
      	// 解密
      	if err = sm4Helper.DecryptCBCFile(encryptedFile, decryptedFile, sem); err != nil {
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
      }
      ```

   3. 加解密大文件（组合*SM2+SM4-CBC*）
      ```go
      // 大文件加密/解密演示（流式 TB级别）
      func TestCBCLargeFileEncryptDecrypt() {
      	var (
      		err           error
      		sem           secret.Semener
      		plainFile           = "/tmp/sm2_test_large_plain.bin"
      		encryptedFile       = "/tmp/sm2_test_large_encrypted.bin"
      		decryptedFile       = "/tmp/sm2_test_large_decrypted.bin"
      		fileSize      int64 = 64 * 1024 * 1024 // 64MB 演示
      		sm4Helper     secret.Symmetricor
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
      	if err = sm4Helper.EncryptCBCLargeFile(plainFile, encryptedFile, sem); err != nil {
      		log.Fatalf("大文件加密失败：%v", err)
      	}
      	log.Printf("大文件加密成功：%s", encryptedFile)
      
      	// 3) 流式解密
      	if err = sm4Helper.DecryptCBCLargeFile(encryptedFile, decryptedFile, sem); err != nil {
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
      }
      
      ```