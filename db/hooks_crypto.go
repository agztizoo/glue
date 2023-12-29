package db

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var (
	EncryptTagName = "encrypt"

	ErrEncryptKeyNotFound = errors.New("encrypt key not found")
)

// CryptoHook 实现数据库加解密能力初始化注入.
//
// 支持的加解密字段类型：string, *string, []byte
//
// ⚠️ 注意: 数据库中的空字符串加载后, 再保存会被加密.
//
// 空串加密:
//   1. 空串加密后存入数据库. 包括: ("", Ptr(""), []byte(""))
//      **兼容不对空串解密做处理的代码**
//   2. 空指针不做处理. 包括: (*string(nil), []byte(nil))
// 空串解密:
//   1. 空字符串不解密.
//   2. 空指针不解密.
// 加解密异常:
//   1. 加密异常报错.
//   2. 解密异常报错.
//
// 对所有标记 encrypt tag 的字段进行加解密处理.
//
// 例:
//	dbopts := config.MySQL
//	aesKeyFunc := util.TenantAesKey
//	cryptoHook := CryptoHook(
//		EncryptTagHandler(aesKeyFunc),
//		DecryptTagHandler(aesKeyFunc),
//	)
//	dial := WithInitializeHook(xxx.Dialector, cryptoHook)
//	dbs := dbopts.OpenDBs(dial, &gorm.Config{})
//  provider := NewProvider(dbs, keyf, scopes...)
func CryptoHook(encrypt, decrypt StringTagHandler) func(*gorm.DB) error {
	p := NewTagProcessor(
		EncryptTagName,
		// 空字符串加密.
		WrapStringTagHandler(encrypt),
		// 空字符串不解密, 不可解密报错.
		WrapStringTagHandler(NonEmptyStringTagHandler(decrypt)))
	return func(db *gorm.DB) error {
		db.Callback().Create().Before("gorm:create").Register("glue:encrypt_field", p.Marshal)
		db.Callback().Create().After("gorm:create").Register("glue:decrypt_field", p.Unmarshal)
		db.Callback().Update().Before("gorm:update").Register("glue:encrypt_field", p.Marshal)
		db.Callback().Update().After("gorm:update").Register("glue:decrypt_field", p.Unmarshal)
		db.Callback().Query().After("gorm:query").Register("glue:decrypt_field", p.Unmarshal)
		return nil
	}
}

// EncryptTagHandler 处理字段解密.
//
// tag    | 算法
// ""       AES
// "true" | AES
// "aes"  | AES
func EncryptTagHandler(keyf func(ctx context.Context) (string, error)) StringTagHandler {
	return func(ctx context.Context, tagValue string, fieldValue string) (string, error) {
		aesKey, err := keyf(ctx)
		if err != nil {
			return "", err
		}
		if aesKey == "" {
			return "", ErrEncryptKeyNotFound
		}

		switch tagValue {
		case "", "true", "aes":
			// 加密
			encrypted, err := AESEncrypt([]byte(fieldValue), []byte(aesKey))
			if err != nil {
				return "", err
			}
			// 将加密后的字符数据用base64编码成字符串
			return base64.StdEncoding.EncodeToString(encrypted), nil
		default:
			return "", fmt.Errorf("encrypt algorithm: %s not support", tagValue)
		}
	}
}

// DecryptTagHandler 处理字段解密.
//
// tag    | 算法
// ""       AES
// "true" | AES
// "aes"  | AES
func DecryptTagHandler(keyf func(ctx context.Context) (string, error)) StringTagHandler {
	return func(ctx context.Context, tagValue string, fieldValue string) (string, error) {
		aesKey, err := keyf(ctx)
		if err != nil {
			return "", err
		}
		if aesKey == "" {
			return "", ErrEncryptKeyNotFound
		}
		switch tagValue {
		case "", "true", "aes":
			// 先用 base64 将字符串解码成字节数组
			encrypted, err := base64.StdEncoding.DecodeString(fieldValue)
			if err != nil {
				return "", err
			}

			// 解密
			decrypted, err := AESDecrypt(encrypted, []byte(aesKey))
			if err != nil {
				return "", err
			}
			return string(decrypted), nil
		default:
			return "", fmt.Errorf("encrypt algorithm: %s not support", tagValue)
		}
	}
}

// AESEncrypt 实现 aes 加密字节数组.
func AESEncrypt(original, aesKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	original = padding(original, block.BlockSize())
	encrypted := make([]byte, len(original))

	blockMode := cipher.NewCBCEncrypter(block, aesKey)
	blockMode.CryptBlocks(encrypted, original)
	return encrypted, nil
}

// AESDecrypt 实现 aes 算法解密字节.
func AESDecrypt(encrypted, aesKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	decrypted := make([]byte, len(encrypted))
	blockMode := cipher.NewCBCDecrypter(block, aesKey)
	blockMode.CryptBlocks(decrypted, encrypted)

	decrypted = unPadding(decrypted)
	return decrypted, nil
}

func padding(src []byte, blockSize int) []byte {
	// 取值1~blockSize，保证padNum被放在串里面了。
	padNum := blockSize - len(src)%blockSize
	pad := bytes.Repeat([]byte{byte(padNum)}, padNum)
	return append(src, pad...)
}

func unPadding(src []byte) []byte {
	n := len(src)
	unPaddingNum := int(src[n-1])
	return src[:n-unPaddingNum]
}
