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
