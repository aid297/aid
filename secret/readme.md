## *Secret* 加密模块

### 目录

| 章节 | 内容 |
|------|------|
| [1. 统一接口说明](#1-统一接口说明) | 选项函数对比、方法列表 |
| [2. 非对称加密（SM2）](#2-非对称加密sm2) | 密钥生成、加解密、签名验证 |
| [3. 非对称加密（RSA）](#3-非对称加密rsa) | 密钥生成、加解密、签名验证 |
| [4. 对称加密（SM4）](#4-对称加密sm4) | ECB/CBC/CTR/GCM 模式，场景选择指南 |
| [5. 对称加密（AES）](#5-对称加密aes) | ECB/CBC/CTR/GCM 模式，AES-192/256 |
| [6. 组合用法](#6-组合用法) | SM2+SM4 / RSA+AES 文件加密 |
| [7. JWT 工具（secret/jwt）](#7-jwt-工具secretjwt) | 令牌工具包：编解码与校验，注入 `secret.Asymmetric` |

1. 统一接口说明

   #### 1.1 非对称加密接口

   **Semen（种子）接口：**

   非对称种子类型为 **`secret.Semen`**；选项属性为 **`secret.SemenAttr`**（如 `rsa.PriKeyBytes`）；公钥/私钥在类型上分别为 **`secret.SemenPubKey`** / **`secret.SemenPriKey`**（与 `crypto` 中公钥私钥概念对应）。

   | 方法 | 说明 |
   |------|------|
   | `NewSem(...)` | 创建种子，不传参则自动生成私钥；传入私钥或公钥参数时验证合法性 |
   | `GeneratePriKey()` | 生成随机私钥，同时验证公私钥合法性 |
   | `GetPriKeyBytes()` / `GetPriKeyBase64()` | 获取私钥 |
   | `GetPubKeyBytes()` / `GetPubKeyBase64()` | 获取公钥 |

   **Asymmetric（加解密/签名）接口：**
   | 方法 | 说明 |
   |------|------|
   | `Encrypt(plainText)` | 加密，返回 base64 密文 |
   | `Decrypt(cipherBase64)` | 解密 |
   | `Sign(data)` | 签名，返回 hex 签名 |
   | `Verify(data, sigHex)` | 验签 |

   #### 1.2 对称加密接口

   对称算法通过 **`secret.Symmetric`** 使用；`sm4.New` / `aes.New` 接受 **`secret.SymmetricAttr`** 选项函数（如 `sm4.KeyBytes`、`aes.AlgorithmCBC` 等）。文件/大文件加解密方法需传入 **`secret.Asymmetric`** 以包装对称密钥。

   **模式选择：**
   | 模式    | 选项函数                                    | Nonce/IV 长度 | 输出格式                   | 适用场景                  |
   | ------- | ------------------------------------------- | ------------- | -------------------------- | ------------------------- |
   | **ECB** | `sm4.AlgorithmECB()` / `aes.AlgorithmECB()` | 不需要        | 密文                       | ❌ 不推荐（不安全）        |
   | **CBC** | `sm4.AlgorithmCBC()` / `aes.AlgorithmCBC()` | 16 字节 IV    | 密文                       | ✅ 通用场景、文件加密      |
   | **CTR** | `sm4.AlgorithmCTR()` / `aes.AlgorithmCTR()` | 16 字节 nonce | nonce(16) + 密文           | ✅ 流式数据、大文件流加密  |
   | **GCM** | `sm4.AlgorithmGCM()` / `aes.AlgorithmGCM()` | 12 字节 nonce | nonce(12) + 密文 + tag(16) | ✅ 认证加密、API 加密、TLS |

   **随机数生成：**
   - `sm4.RandKey()` / `aes.RandKey()` - 生成随机密钥
   - `sm4.RandIV()` / `aes.RandIV()` - 生成随机 IV（16 字节）
   - `sm4.RandCTRNonce()` / `aes.RandCTRNonce()` - 生成 16 字节随机 nonce（CTR 模式）
   - `sm4.RandGCMNonce()` / `aes.RandGCMNonce()` - 生成 12 字节随机 nonce（GCM 模式）

   **默认模式为 CBC**

   **统一方法列表：**
   | 方法 | 说明 |
   |------|------|
   | `Encrypt(plainText)` | 加密，返回密文 |
   | `Decrypt(cipherText)` | 解密 |
   | `EncryptBase64/DecryptBase64` | Base64 编解码 |
   | `EncryptStream/DecryptStream` | 流式加解密 |
   | `EncryptFile/DecryptFile` | 文件加解密 |
   | `EncryptLargeFile/DecryptLargeFile` | 大文件流式加解密 |

2. 非对称加密（*SM2*）

   1. 生成种子
      ```go
      func TestGenerateKeyPair(t *testing.T) {
      	var (
      		err                        error
      		sem                        secret.Semen
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
      		semEncrypter, semDecrypter secret.Semen
      		sm2Encrypter, sm2Decrypter secret.Asymmetric
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
      		semSign, semVerify secret.Semen
      		sm2Sign, sm2Verify secret.Asymmetric
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
      		semSign, semVerify secret.Semen
      		sm2Sign, sm2Verify secret.Asymmetric
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

3. 非对称加密（*RSA*）

   1. 生成种子
      ```go
      func TestGenerateKeyPair(t *testing.T) {
      	var (
      		err                        error
      		sem                        secret.Semen
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
      		semA, semB   secret.Semen
      		rsaA, rsaB   secret.Asymmetric
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
      		semSign, semVerify secret.Semen
      		rsaSign, rsaVerify secret.Asymmetric
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
      		semSign, semVerify secret.Semen
      		rsaSign, rsaVerify secret.Asymmetric
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

4. 对称加密（*SM4*）

   *SM4* 支持 *ECB*、*CBC*、*CTR*、*GCM* 四种模式，通过 `AlgorithmCBC()`/`AlgorithmECB()`/`AlgorithmCTR()`/`AlgorithmGCM()` 选择。默认 *CBC* 模式。

   ### 模式适用场景速查

   | 模式 | 特性 | 适用场景 | 不适用场景 |
   |------|------|----------|------------|
   | **ECB** | 各区块独立加密，相同明文产生相同密文 | 学习/测试，理论演示 | **敏感数据加密**（会暴露明文模式） |
   | **CBC** | 需 IV，区块链式依赖，支持填充 | 文件加密、数据库字段、通用数据 | 流式加密、大文件流处理（需填充） |
   | **CTR** | 流式加密，无需填充，16字节 nonce | 网络通信、流式数据、大文件流加密 | 需要完整性认证的场景 |
   | **GCM** | 认证加密，12字节 nonce，自动附加 16 字节 tag | TLS/SSL、API 加密、认证通信、敏感数据传输 | 简单快速加密（开销比 CTR 大） |

   **推荐原则：**
   - 需要认证：`GCM`
   - 流式加密无需认证：`CTR`
   - 通用文件加密：`CBC` 或 `GCM`
   - **永远不要用 ECB 加密真实敏感数据**

   ### 各模式详解

   #### 3.1 CBC 模式（默认，推荐通用场景）

   CBC（Cipher Block Chaining）模式通过将前一个密文块与当前明文块异或后再加密，安全性高，适用于大多数场景。

   **基本用法（默认 *CBC* 模式）
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

   #### 3.2 ECB 模式（不推荐，仅供学习）

   ⚠️ **警告**：ECB 模式下相同的明文块总是产生相同的密文块，会暴露明文的结构模式，**请勿用于加密敏感数据**。

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

   #### 3.3 Base64 编解码
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

   #### 3.4 CTR 模式（流式加密，无需填充）

   CTR（Counter）模式将分组密码转换为流密码，无需填充，适合流式数据加密和大文件流处理。

   - **Nonce 长度**：16 字节（与 blockSize 相同）
   - **输出格式**：`nonce(16) + 密文`
   - **特点**：无需填充，可并行加密，速度快

   ```go
   func TestCTR(t *testing.T) {
       var (
           testKey   = []byte("1234567890abcdef") // 16 字节密钥
           testNonce = []byte("abcdef1234567890") // 16 字节 nonce
           testPlain = []byte("hello, SM4 CTR encrypt test!")
       )
   
       sm4Helper, err := sm4.New(sm4.KeyBytes(testKey), sm4.IVBytes(testNonce), sm4.AlgorithmCTR())
       if err != nil {
           t.Fatalf("创建 CTR 对象失败：%v", err)
       }
       cipherText, err := sm4Helper.Encrypt(testPlain)
       if err != nil {
           t.Fatalf("加密失败：%v", err)
       }
       plain, err := sm4Helper.Decrypt(cipherText)
       if err != nil {
           t.Fatalf("解密失败：%v", err)
       }
       if !bytes.Equal(plain, testPlain) {
           t.Fatalf("比对失败：结果 %s，期望 %s", plain, testPlain)
       }
       t.Logf("CTR OK: %s", plain)
   }
   
   // 使用随机 nonce
   func TestCTRWithRandNonce(t *testing.T) {
       var testKey = []byte("1234567890abcdef")
       var testPlain = []byte("hello, SM4 CTR with random nonce!")
   
       sm4Helper, err := sm4.New(sm4.KeyBytes(testKey), sm4.RandCTRNonce(), sm4.AlgorithmCTR())
       if err != nil {
           t.Fatalf("创建 CTR 对象失败：%v", err)
       }
       cipherText, err := sm4Helper.Encrypt(testPlain)
       if err != nil {
           t.Fatalf("加密失败：%v", err)
       }
       plain, err := sm4Helper.Decrypt(cipherText)
       if err != nil {
           t.Fatalf("解密失败：%v", err)
       }
       if !bytes.Equal(plain, testPlain) {
           t.Fatalf("比对失败：结果 %s，期望 %s", plain, testPlain)
       }
       t.Logf("CTR RandNonce OK: %s", plain)
   }
   ```

   #### 3.5 GCM 模式（认证加密）

   GCM（Galois/Counter Mode）是认证加密模式，同时提供机密性和完整性认证。解密时会自动验证数据完整性。

   - **Nonce 长度**：12 字节（标准）
   - **输出格式**：`nonce(12) + 密文 + tag(16)`
   - **特点**：自动认证，防篡改，解密失败返回错误

   ```go
   func TestGCM(t *testing.T) {
       var (
           testKey   = []byte("1234567890abcdef") // 16 字节密钥
           testNonce = []byte("abcdef123456")      // 12 字节 nonce
           testPlain = []byte("hello, SM4 GCM encrypt test!")
       )
   
       sm4Helper, err := sm4.New(sm4.KeyBytes(testKey), sm4.IVBytes(testNonce), sm4.AlgorithmGCM())
       if err != nil {
           t.Fatalf("创建 GCM 对象失败：%v", err)
       }
       cipherText, err := sm4Helper.Encrypt(testPlain)
       if err != nil {
           t.Fatalf("加密失败：%v", err)
       }
       t.Logf("GCM 密文长度：%d 字节", len(cipherText)) // 52 = 12(nonce) + 24(加密) + 16(tag)
   
       plain, err := sm4Helper.Decrypt(cipherText)
       if err != nil {
           t.Fatalf("解密失败：%v", err)
       }
       if !bytes.Equal(plain, testPlain) {
           t.Fatalf("比对失败：结果 %s，期望 %s", plain, testPlain)
       }
       t.Logf("GCM OK: %s", plain)
   }
   
   // 篡改检测
   func TestGCMTampered(t *testing.T) {
       testNonce := []byte("abcdef123456")
       sm4Helper, err := sm4.New(sm4.KeyBytes(testKey), sm4.IVBytes(testNonce), sm4.AlgorithmGCM())
       if err != nil {
           t.Fatalf("创建 GCM 对象失败：%v", err)
       }
       cipherText, err := sm4Helper.Encrypt(testPlain)
       if err != nil {
           t.Fatalf("加密失败：%v", err)
       }
   
       // 篡改密文
       if len(cipherText) > 20 {
           cipherText[20] ^= 0xFF
       }
   
       _, err = sm4Helper.Decrypt(cipherText)
       if err == nil {
           t.Fatal("GCM 认证应该失败：密文已被篡改")
       }
       t.Logf("GCM 篡改检测成功：%v", err)
   }
   
   // 使用随机 nonce
   func TestGCMWithRandNonce(t *testing.T) {
       var testKey = []byte("1234567890abcdef")
       var testPlain = []byte("hello, SM4 GCM with random nonce!")
   
       sm4Helper, err := sm4.New(sm4.KeyBytes(testKey), sm4.RandGCMNonce(), sm4.AlgorithmGCM())
       if err != nil {
           t.Fatalf("创建 GCM 对象失败：%v", err)
       }
       cipherText, err := sm4Helper.Encrypt(testPlain)
       if err != nil {
           t.Fatalf("加密失败：%v", err)
       }
       plain, err := sm4Helper.Decrypt(cipherText)
       if err != nil {
           t.Fatalf("解密失败：%v", err)
       }
       if !bytes.Equal(plain, testPlain) {
           t.Fatalf("比对失败：结果 %s，期望 %s", plain, testPlain)
       }
       t.Logf("GCM RandNonce OK: %s", plain)
   }
   ```

5. 对称加密（*AES*）

   AES 支持 ECB、CBC、CTR、GCM 四种模式，与 SM4 类似。默认 CBC 模式。支持 AES-128、AES-192、AES-256 三种密钥长度。

   #### 4.1 CBC 模式（默认，推荐通用场景）

   **基本用法（默认 CBC 模式）**
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

   #### 4.2 ECB 模式（不推荐，仅供学习）

   ⚠️ **警告**：ECB 模式下相同的明文块总是产生相同的密文块，**请勿用于加密敏感数据**。
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

   #### 4.3 AES192 和 AES256 支持
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

   #### 4.4 CTR 模式（流式加密，无需填充）

   CTR（Counter）模式将分组密码转换为流密码，无需填充，适合流式数据加密。

   - **Nonce 长度**：16 字节（与 blockSize 相同）
   - **输出格式**：`nonce(16) + 密文`

   ```go
   func TestCTR(t *testing.T) {
       var (
           testKey   = []byte("1234567890abcdef") // 16 字节密钥
           testNonce = []byte("abcdef1234567890") // 16 字节 nonce
           testPlain = []byte("hello, AES CTR encrypt test!")
       )
   
       aesHelper, err := aes.New(aes.KeyBytes(testKey), aes.IVBytes(testNonce), aes.AlgorithmCTR())
       if err != nil {
           t.Fatalf("创建 CTR 对象失败：%v", err)
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
       t.Logf("CTR OK: %s", plain)
   }
   
   func TestCTRBase64(t *testing.T) {
       var (
           testKey   = []byte("1234567890abcdef")
           testNonce = []byte("abcdef1234567890")
           testPlain = []byte("hello, AES CTR encrypt test!")
       )
   
       aesHelper, err := aes.New(aes.KeyBytes(testKey), aes.IVBytes(testNonce), aes.AlgorithmCTR())
       if err != nil {
           t.Fatalf("创建 CTR 对象失败：%v", err)
       }
       b64, err := aesHelper.EncryptBase64(testPlain)
       if err != nil {
           t.Fatalf("加密失败：%v", err)
       }
       t.Logf("加密后 base64：%s", b64)
       plain, err := aesHelper.DecryptBase64(b64)
       if err != nil {
           t.Fatalf("解密失败：%v", err)
       }
       if !bytes.Equal(plain, testPlain) {
           t.Fatalf("比对失败：结果 %s，期望 %s", plain, testPlain)
       }
       t.Logf("CTR Base64 OK: %s", plain)
   }
   
   // 使用随机 nonce
   func TestCTRWithRandNonce(t *testing.T) {
       var (
           testKey   = []byte("1234567890abcdef")
           testPlain = []byte("hello, AES CTR with random nonce!")
       )
   
       aesHelper, err := aes.New(aes.KeyBytes(testKey), aes.RandCTRNonce(), aes.AlgorithmCTR())
       if err != nil {
           t.Fatalf("创建 CTR 对象失败：%v", err)
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
       t.Logf("CTR RandNonce OK: %s", plain)
   }
   ```

   #### 4.5 GCM 模式（认证加密）

   GCM（Galois/Counter Mode）是认证加密模式，同时提供机密性和完整性认证。

   - **Nonce 长度**：12 字节（标准）
   - **输出格式**：`nonce(12) + 密文 + tag(16)`

   ```go
   func TestGCM(t *testing.T) {
       var (
           testKey   = []byte("1234567890abcdef")
           testNonce = []byte("abcdef123456")      // 12 字节 nonce
           testPlain = []byte("hello, AES GCM encrypt test!")
       )
   
       aesHelper, err := aes.New(aes.KeyBytes(testKey), aes.IVBytes(testNonce), aes.AlgorithmGCM())
       if err != nil {
           t.Fatalf("创建 GCM 对象失败：%v", err)
       }
       cipherText, err := aesHelper.Encrypt(testPlain)
       if err != nil {
           t.Fatalf("加密失败：%v", err)
       }
       t.Logf("GCM 密文长度：%d 字节", len(cipherText)) // 56 = 12(nonce) + 28(加密) + 16(tag)
   
       plain, err := aesHelper.Decrypt(cipherText)
       if err != nil {
           t.Fatalf("解密失败：%v", err)
       }
       if !bytes.Equal(plain, testPlain) {
           t.Fatalf("比对失败：结果 %s，期望 %s", plain, testPlain)
       }
       t.Logf("GCM OK: %s", plain)
   }
   
   func TestGCMBase64(t *testing.T) {
       var (
           testKey   = []byte("1234567890abcdef")
           testNonce = []byte("abcdef123456")
           testPlain = []byte("hello, AES GCM encrypt test!")
       )
   
       aesHelper, err := aes.New(aes.KeyBytes(testKey), aes.IVBytes(testNonce), aes.AlgorithmGCM())
       if err != nil {
           t.Fatalf("创建 GCM 对象失败：%v", err)
       }
       b64, err := aesHelper.EncryptBase64(testPlain)
       if err != nil {
           t.Fatalf("加密失败：%v", err)
       }
       t.Logf("加密后 base64：%s", b64)
       plain, err := aesHelper.DecryptBase64(b64)
       if err != nil {
           t.Fatalf("解密失败：%v", err)
       }
       if !bytes.Equal(plain, testPlain) {
           t.Fatalf("比对失败：结果 %s，期望 %s", plain, testPlain)
       }
       t.Logf("GCM Base64 OK: %s", plain)
   }
   
   // 使用随机 nonce
   func TestGCMWithRandNonce(t *testing.T) {
       var (
           testKey   = []byte("1234567890abcdef")
           testPlain = []byte("hello, AES GCM with random nonce!")
       )
   
       aesHelper, err := aes.New(aes.KeyBytes(testKey), aes.RandGCMNonce(), aes.AlgorithmGCM())
       if err != nil {
           t.Fatalf("创建 GCM 对象失败：%v", err)
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
       t.Logf("GCM RandNonce OK: %s", plain)
   }
   
   // 篡改检测
   func TestGCMTampered(t *testing.T) {
       var (
           testKey   = []byte("1234567890abcdef")
           testNonce = []byte("abcdef123456")
           testPlain = []byte("hello, AES GCM tamper test!")
       )
   
       aesHelper, err := aes.New(aes.KeyBytes(testKey), aes.IVBytes(testNonce), aes.AlgorithmGCM())
       if err != nil {
           t.Fatalf("创建 GCM 对象失败：%v", err)
       }
       cipherText, err := aesHelper.Encrypt(testPlain)
       if err != nil {
           t.Fatalf("加密失败：%v", err)
       }
   
       // 篡改密文
       if len(cipherText) > 20 {
           cipherText[20] ^= 0xFF
       }
   
       _, err = aesHelper.Decrypt(cipherText)
       if err == nil {
           t.Fatal("GCM 认证应该失败：密文已被篡改")
       }
       t.Logf("GCM 篡改检测成功：%v", err)
   }
   ```

6. 组合用法

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
      		semEncrypter, semDecrypter secret.Semen
      		sm2Encrypter, sm2Decrypter secret.Asymmetric
      		semPriKeyEncrypter         []byte
      		plainFile                  = "/tmp/sm2_test_plain.txt"
      		encryptedFile              = "/tmp/sm2_test_encrypted.bin"
      		decryptedFile              = "/tmp/sm2_test_decrypted.txt"
      		sm4Encrypter, sm4Decrypter secret.Symmetric
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
      		semEncrypter, semDecrypter secret.Semen
      		sm2Encrypter, sm2Decrypter secret.Asymmetric
      		sm4Encrypter, sm4Decrypter secret.Symmetric
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
      		semEncrypter, semDecrypter secret.Semen
      		rsaEncrypter, rsaDecrypter secret.Asymmetric
      		aesEncrypter, aesDecrypter secret.Symmetric
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
      		semEncrypter, semDecrypter secret.Semen
      		rsaEncrypter, rsaDecrypter secret.Asymmetric
      		aesEncrypter, aesDecrypter secret.Symmetric
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

7. JWT 工具（*secret/jwt*）

   本节放在文档末尾：JWT 属于 **secret 下的令牌工具包**（包路径 `github.com/aid297/aid/v2/secret/jwt`），与 `secret/asymmetric` 中的 RSA、ECDSA、Ed25519、SM2 等**并列**，通过注入 `secret.Asymmetric` 完成签名与验签，而不作为某一类非对称子算法目录的一部分。

   JWT（JSON Web Token）用于在各方之间安全传输声明。

   ### 声明结构

   ```go
   type Claims struct {
       Iss   string                 // 签发者
       Sub   string                 // 主题
       Aud   string                 // 受众
       Exp   int64                  // 过期时间（Unix 时间戳）
       Nbf   int64                  // 生效时间（Unix 时间戳）
       Iat   int64                  // 签发时间（Unix 时间戳）
       Jti   string                 // JWT ID
       Extra map[string]interface{}  // 自定义声明
   }
   ```

   ### 基本用法

   ```go
   import (
       "testing"
       "time"

       "github.com/aid297/aid/v2/secret"
       "github.com/aid297/aid/v2/secret/asymmetric/rsa"
       "github.com/aid297/aid/v2/secret/jwt"
   )

   func TestJWTBasic(t *testing.T) {
       // 1. 创建 RSA 密钥对
       rsaSem, err := rsa.NewSem()
       if err != nil {
           t.Fatalf("生成密钥对失败: %v", err)
       }

       // 2. 创建 JWT 实例（使用 secret.Asymmetric）
       var asymm secret.Asymmetric = rsa.New(rsaSem)
       jwtInstance := jwt.New(asymm)

       // 3. 构建声明
       claims := &jwt.Claims{
           Iss: "test-issuer",
           Sub: "test-subject",
           Aud: "test-audience",
           Iat: time.Now().Unix(),
           Exp: time.Now().Add(time.Hour).Unix(),
           Nbf: time.Now().Unix() - 60,
           Jti: "unique-token-id",
       }

       // 4. 生成 token
       token, err := jwtInstance.Generate(claims)
       if err != nil {
           t.Fatalf("生成 JWT 失败: %v", err)
       }
       t.Logf("Token: %s", token)

       // 5. 验证 token
       verifiedClaims, err := jwtInstance.Verify(token)
       if err != nil {
           t.Fatalf("验证 JWT 失败: %v", err)
       }
       t.Logf("验证通过: %+v", verifiedClaims)
   }
   ```

   ### 使用现有密钥对

   ```go
   import (
       "testing"
       "time"

       "github.com/aid297/aid/v2/secret"
       "github.com/aid297/aid/v2/secret/asymmetric/rsa"
       "github.com/aid297/aid/v2/secret/jwt"
   )

   func TestJWTWithExistingKeys(t *testing.T) {
       // 生成密钥对
       rsaSem, err := rsa.NewSem()
       if err != nil {
           t.Fatalf("生成密钥对失败: %v", err)
       }

       // 获取公私钥
       pubKeyBytes, _ := rsaSem.GetPubKeyBytes()
       priKeyBytes, _ := rsaSem.GetPriKeyBytes()

       // 使用公钥创建 JWT（仅验证）
       rsaSemPub, _ := rsa.NewSem(rsa.PubKeyBytes(pubKeyBytes))
       // 使用私钥创建 JWT（仅签名）
       rsaSemPri, _ := rsa.NewSem(rsa.PriKeyBytes(priKeyBytes))

       // 签名
       var signerAsymm secret.Asymmetric = rsa.New(rsaSemPri)
       token, _ := jwt.New(signerAsymm).Generate(&jwt.Claims{
           Iss: "test",
           Exp: time.Now().Add(time.Hour).Unix(),
       })

       // 验证
       var verifierAsymm secret.Asymmetric = rsa.New(rsaSemPub)
       claims, _ := jwt.New(verifierAsymm).Verify(token)
       t.Logf("Iss: %s", claims.Iss)
   }
   ```

   ### 支持多算法

   通过 `secret.Asymmetric` 注入算法；头部 `alg` 使用 **`jwt.NewWithAlg(alg, asymm)`**（`alg` 类型为 `jwt.Alg`）。不写 `alg` 时用 **`jwt.New(asymm)`**。

   - **RS256** / **RS384** / **RS512**：RSA + PKCS#1 v1.5 + 对应 SHA（`rsa.New`）  
   - **ES256** / **ES384** / **ES512**：ECDSA + 对应 SHA（`ecdsa.New`，当前种子默认 P-256 对应 **ES256**）  
   - **EdDSA**：Ed25519（`ed25519.New`）  
   - **SM2**：国密 SM2（`sm2.New`）  

   只需替换 `secret.Asymmetric` 实现及 `jwt.Alg` 即可切换算法。