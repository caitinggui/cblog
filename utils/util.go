package utils

import (
	"bytes"
	"context"
	"crypto/des"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/time/rate"

	logger "github.com/caitinggui/seelog"
	"github.com/caitinggui/uniqueid"
)

var UID *uniqueid.UniqueId
var HttpClient = http.Client{Timeout: 3 * time.Second}

func IfPathExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func StrToUint(s string) uint {
	n, _ := strconv.ParseUint(s, 10, 64)
	return uint(n)
}

func StrToUint64(s string) (n uint64) {
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

// 将任意类型(包括时间)转字符串
func ToStr(value interface{}) (s string) {
	switch v := value.(type) {
	case bool:
		s = strconv.FormatBool(v)
	case float32:
		s = strconv.FormatFloat(float64(v), 'f', 3, 32)
	case float64:
		s = strconv.FormatFloat(v, 'f', 3, 64)
	case int:
		s = strconv.FormatInt(int64(v), 10)
	case int8:
		s = strconv.FormatInt(int64(v), 10)
	case int16:
		s = strconv.FormatInt(int64(v), 10)
	case int32:
		s = strconv.FormatInt(int64(v), 10)
	case int64:
		s = strconv.FormatInt(v, 10)
	case uint:
		s = strconv.FormatUint(uint64(v), 10)
	case uint8:
		s = strconv.FormatUint(uint64(v), 10)
	case uint16:
		s = strconv.FormatUint(uint64(v), 10)
	case uint32:
		s = strconv.FormatUint(uint64(v), 10)
	case uint64:
		s = strconv.FormatUint(v, 10)
	case string:
		s = v
	case []byte:
		s = string(v)
	case time.Time:
		s = v.Format("2006-01-02 15:04:05")
	default:
		s = fmt.Sprintf("%v", v)
	}
	return s
}

func InitUniqueId(WorkerId uint16, ReserveId uint8) {
	UID = uniqueid.NewUniqueId(WorkerId, ReserveId)
}

// 生成唯一id，要先调用InitUniqueId
func GenerateId() uint64 {
	uid, err := UID.NextId()
	if err != nil {
		logger.Error("生成uid失败")
		panic(err)
	}
	return uid
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
	out, err := DesEcbDecrypt(data, []byte("rysbkeug"))
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

type JsonTime time.Time

func (self JsonTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(self).Format("2006-01-02T15:04:05.999"))
	return []byte(stamp), nil
}

// 3次重试
func HttpRetryGet(url string) (body []byte, err error) {
	var (
		resp          *http.Response
		maxRetryTimes = 3 // 3次重试
	)

	for i := 0; i < maxRetryTimes; i++ {
		logger.Info(url, " retryTimes: ", i)
		resp, err = HttpClient.Get(url)
		// err == nil时才需要关闭body
		if err == nil {
			body, _ = ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				break
			}
		}
	}
	return
}

func HttpGet(url string) (body []byte, err error) {
	var resp *http.Response
	resp, err = HttpClient.Get(url)
	if err != nil {
		logger.Error("request get error: ", url, err)
		return body, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return
}

func HttpPost(url string, data *map[string]interface{}) (body []byte, err error) {
	var resp *http.Response
	var jsonValue []byte
	contentType := "application/json"
	jsonValue, err = json.Marshal(data)
	if err != nil {
		logger.Warn("request parameter error: ", err)
		return
	}

	resp, err = HttpClient.Post(url, contentType, bytes.NewBuffer(jsonValue))
	if err != nil {
		logger.Error("request post error: ", url, err)
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return
}

func PostJsonWithAuth(url string, values *map[string]interface{}, accessToken string) (body []byte, err error) {
	contentType := "application/json"
	jsonValue, _ := json.Marshal(values)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", accessToken)

	resp, err1 := HttpClient.Do(req)
	if err1 != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return
}

type RateLimiter struct {
	Capacity int
	Interval time.Duration
	Limiter  *rate.Limiter
	Ctx      context.Context
}

func (self RateLimiter) Wait() error {
	return self.Limiter.Wait(self.Ctx)
}

func NewRateLimter(Interval time.Duration, Capacity int) *RateLimiter {
	ctx := context.Background()
	rlm := RateLimiter{
		Capacity: Capacity,
		Interval: Interval,
		Ctx:      ctx,
		Limiter:  rate.NewLimiter(rate.Every(Interval), Capacity),
	}
	return &rlm
}

// 定时任务，指定每天几点几分几秒执行任务
func CronFuncDaily(f func() error, hour, min, sec int) {
	timer := time.NewTimer(time.Second)
	for {
		logger.Info("开始执行定时任务: ", f)
		err := f()
		now := time.Now()
		nextDay := now.Add(time.Hour * 24)
		nextDay = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), hour, min, sec, 0, nextDay.Location())
		logger.Info("定时任务执行完成: ", err, " 开始下一次定时: ", nextDay)
		timer.Reset(nextDay.Sub(now))
		<-timer.C
	}
}

func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func FormatIP(ip string) string {
	res := strings.Split(ip, ".")
	if len(res) < 2 {
		return res[0] + "*"
	}
	return res[0] + ".*.*." + res[len(res)-1]
}
