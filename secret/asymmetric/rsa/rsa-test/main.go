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

	if aesEncrypter, err = aes.New(aes.RandKeyWithBits(aes.AES192, &aesKey), aes.RandIV(&aesIV)); err != nil {
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

	if aesEncrypter, err = aes.New(aes.RandKeyWithBits(aes.AES256), aes.RandIV()); err != nil {
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
