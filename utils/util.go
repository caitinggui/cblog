package utils

import (
	"bytes"
	"crypto/des"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"

	logger "github.com/caitinggui/seelog"
)

func StrToUnit64(s string) (n uint64) {
	// 10进制，64位
	n, _ = strconv.ParseUint(s, 10, 64)
	return n
}

func StrToInt64(s string) (n int64) {
	n, _ = strconv.ParseInt(s, 10, 64)
	return
}

func StrToFloat64(s string) (n float64) {
	n, _ = strconv.ParseFloat(s, 64)
	return
}

// 生成人类易读随机字符串, 不包含0,1,l,o,I,O等字符
func RandomHumanString(length int) string {
	// 去除0,1,l,o,I,O等难以区分的字
	str := "23456789abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
	byteStr := []byte(str)
	result := make([]byte, 0, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, byteStr[r.Intn(len(byteStr))])
	}
	return string(result)
}

// 解密uid信息
func UidDecrypt(hexStr string) (outString string) {
	outString = ""
	hexStr = strings.Replace(hexStr, "-", "", -1)
	data, _ := hex.DecodeString(hexStr)
	out, err := DesEcbDecrypt(data, []byte("yanglong"))
	if err != nil {
		return
	}
	outString = string(out)
	return outString
}

// 加密uid信息,返回十六进制的字符串
func UidEncrypt(inString string) (string, error) {
	out, err := DesEcbEncrypt([]byte(inString), []byte("yanglong"))
	logger.Info("out: ", out)
	hexStr := hex.EncodeToString(out)
	return hexStr, err
}

// des填充方式
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// des填充方式
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func DesEcbEncrypt(data, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	data = PKCS5Padding(data, bs)
	if len(data)%bs != 0 {
		return nil, errors.New("Need a multiple of the blocksize")
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	return out, nil
}

func DesEcbDecrypt(data []byte, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	if len(data)%bs != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		block.Decrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	out = PKCS5UnPadding(out)
	return out, nil
}

func HmacSha256(data, key string) string {
	hm := hmac.New(sha256.New, []byte(key))
	hm.Write([]byte(data))
	return hex.EncodeToString(hm.Sum(nil))
}
