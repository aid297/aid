package main

import (
	"crypto/rand"
	"errors"
	"io"
	"log"
	"os"

	"github.com/aid297/aid/v3/secrets"
	"github.com/aid297/aid/v3/secrets/asymmetric/sm2"
	"github.com/aid297/aid/v3/secrets/symmetric/sm4"
)

// TestFileEncrypt 1. 文件加密/解密演示
func TestFileEncrypt() {
	// 生成密钥对
	var (
		err                        error
		semEncrypter, semDecrypter secrets.Semen
		sm2Encrypter, sm2Decrypter secrets.Asymmetric
		semPriKeyEncrypter         []byte
		plainFile                  = "/tmp/sm2_test_plain.txt"
		encryptedFile              = "/tmp/sm2_test_encrypted.bin"
		decryptedFile              = "/tmp/sm2_test_decrypted.txt"
		sm4Encrypter, sm4Decrypter secrets.Symmetric
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
		semEncrypter, semDecrypter secrets.Semen
		sm2Encrypter, sm2Decrypter secrets.Asymmetric
		sm4Encrypter, sm4Decrypter secrets.Symmetric
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
