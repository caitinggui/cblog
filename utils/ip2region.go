package utils

import (
	//"io/ioutil"
	//"net/http"
	//"time"

	"errors"
	logger "github.com/caitinggui/seelog"
	"github.com/json-iterator/go"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var dataPath = "./static/libs/ip2region/ip2region.db"
var onceIp2Region sync.Once
var ip2Region *Ip2Region
var ipInfo IpInfo

func getIp2Region() *Ip2Region {
	onceIp2Region.Do(func() {
		logger.Info("check if ip2region dataPath exist: ", dataPath)
		if !IfPathExist(dataPath) {
			return
		}
		logger.Info(dataPath, " exist")
		ip2Region, _ = NewIp2Region(dataPath)
	})
	return ip2Region
}

type IpInfo struct {
	IP       string `json:"ip"`
	Country  string `json:"country"`
	Province string `json:"region"`
	City     string `json:"city"`
	ISP      string `json:"isp"`
	//Timeout  time.Duration // 请求的超时时间
}

// 用淘宝的接口去解析
func (self *IpInfo) Name() string {
	return "淘宝"
}

func (self *IpInfo) Url() string {
	return "http://ip.taobao.com/service/getIpInfo.php?accessKey=alibaba-inc&ip="
}

func (self *IpInfo) PraseIp() error {
	ipInfo, err := PraseIp(self.IP)
	if err != nil {
		return err
	}
	self.Country = ipInfo.Country
	self.City = ipInfo.City
	self.ISP = ipInfo.ISP
	self.Province = ipInfo.Province
	return nil
}

func (self *IpInfo) String() string {
	return self.Country + "|" + self.Province + "|" + self.City + "|" + self.ISP
}

// 先用离线解析，如果离线解析失败或者省份为空，就用网络解析，如果网络解析，就用前面离线解析的结果
func PraseIp(IP string) (ipInfo IpInfo, err error) {
	ipInfo.IP = IP
	if IP == "::1" || IP == "127.0.0.1" {
		ipInfo.ISP = "内网IP"
		return ipInfo, nil
	}
	if ip2Region = getIp2Region(); ip2Region != nil {
		ipInfo, err = ip2Region.MemorySearch(IP)
		if err == nil && ipInfo.Province != "" {
			logger.Info("ip2region res: ", ipInfo.Country, ipInfo.Province)
			return ipInfo, nil
		}
		logger.Warnf("ip2region %p prase ip failed: %v", ip2Region, err)
	}
	url := ipInfo.Url() + IP
	body, err := HttpRetryGet(url)
	if err != nil {
		logger.Error("ip2Region request ip error: ", err)
		return ipInfo, err
	}
	ipInfoFromNet := IpInfo{}
	json.Unmarshal([]byte(jsoniter.Get(body, "data").ToString()), &ipInfoFromNet)
	logger.Infof("%s res: %s %s", ipInfo.Name(), ipInfo.Country, ipInfo.Province)
	if ipInfoFromNet.Country == "" {
		return ipInfo, nil
	}
	ipInfoFromNet.IP = IP
	return ipInfoFromNet, nil
}

const (
	INDEX_BLOCK_LENGTH  = 12
	TOTAL_HEADER_LENGTH = 8192
)

type Ip2Region struct {
	// db file handler
	dbFileHandler *os.File

	//header block info

	headerSip []int64
	headerPtr []int64
	headerLen int64

	// super block index info
	firstIndexPtr int64
	lastIndexPtr  int64
	totalBlocks   int64

	// for memory mode only
	// the original db binary string

	dbBinStr []byte
	dbFile   string
}

func getIpInfo(cityId int64, line []byte) IpInfo {

	lineSlice := strings.Split(string(line), "|")
	ipInfo := IpInfo{}
	length := len(lineSlice)
	if length < 5 {
		for i := 0; i <= 5-length; i++ {
			lineSlice = append(lineSlice, "")
		}
	}

	ipInfo.Country = lineSlice[0]
	ipInfo.Province = lineSlice[2]
	ipInfo.City = lineSlice[3]
	ipInfo.ISP = lineSlice[4]
	return ipInfo
}

func NewIp2Region(path string) (*Ip2Region, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return &Ip2Region{
		dbFile:        path,
		dbFileHandler: file,
	}, nil
}

func (this *Ip2Region) Close() {
	this.dbFileHandler.Close()
}

func (this *Ip2Region) MemorySearch(ipStr string) (ipInfo IpInfo, err error) {
	ipInfo = IpInfo{IP: ipStr}

	if this.totalBlocks == 0 {
		this.dbBinStr, err = ioutil.ReadFile(this.dbFile)

		if err != nil {

			return ipInfo, err
		}

		this.firstIndexPtr = getLong(this.dbBinStr, 0)
		this.lastIndexPtr = getLong(this.dbBinStr, 4)
		this.totalBlocks = (this.lastIndexPtr-this.firstIndexPtr)/INDEX_BLOCK_LENGTH + 1
	}

	ip, err := ip2long(ipStr)
	if err != nil {
		return ipInfo, err
	}

	h := this.totalBlocks
	var dataPtr, l int64
	for l <= h {

		m := (l + h) >> 1
		p := this.firstIndexPtr + m*INDEX_BLOCK_LENGTH
		sip := getLong(this.dbBinStr, p)
		if ip < sip {
			h = m - 1
		} else {
			eip := getLong(this.dbBinStr, p+4)
			if ip > eip {
				l = m + 1
			} else {
				dataPtr = getLong(this.dbBinStr, p+8)
				break
			}
		}
	}
	if dataPtr == 0 {
		return ipInfo, errors.New("not found")
	}

	dataLen := ((dataPtr >> 24) & 0xFF)
	dataPtr = (dataPtr & 0x00FFFFFF)
	ipInfo = getIpInfo(getLong(this.dbBinStr, dataPtr), this.dbBinStr[(dataPtr)+4:dataPtr+dataLen])
	ipInfo.IP = ipStr
	return ipInfo, nil

}

func (this *Ip2Region) BinarySearch(ipStr string) (ipInfo IpInfo, err error) {
	ipInfo = IpInfo{IP: ipStr}
	if this.totalBlocks == 0 {
		this.dbFileHandler.Seek(0, 0)
		superBlock := make([]byte, 8)
		this.dbFileHandler.Read(superBlock)
		this.firstIndexPtr = getLong(superBlock, 0)
		this.lastIndexPtr = getLong(superBlock, 4)
		this.totalBlocks = (this.lastIndexPtr-this.firstIndexPtr)/INDEX_BLOCK_LENGTH + 1
	}

	var l, dataPtr, p int64

	h := this.totalBlocks

	ip, err := ip2long(ipStr)

	if err != nil {
		return
	}

	for l <= h {
		m := (l + h) >> 1

		p = m * INDEX_BLOCK_LENGTH

		_, err = this.dbFileHandler.Seek(this.firstIndexPtr+p, 0)
		if err != nil {
			return
		}

		buffer := make([]byte, INDEX_BLOCK_LENGTH)
		_, err = this.dbFileHandler.Read(buffer)

		if err != nil {

		}
		sip := getLong(buffer, 0)
		if ip < sip {
			h = m - 1
		} else {
			eip := getLong(buffer, 4)
			if ip > eip {
				l = m + 1
			} else {
				dataPtr = getLong(buffer, 8)
				break
			}
		}

	}

	if dataPtr == 0 {
		err = errors.New("not found")
		return
	}

	dataLen := ((dataPtr >> 24) & 0xFF)
	dataPtr = (dataPtr & 0x00FFFFFF)

	this.dbFileHandler.Seek(dataPtr, 0)
	data := make([]byte, dataLen)
	this.dbFileHandler.Read(data)
	ipInfo = getIpInfo(getLong(data, 0), data[4:dataLen])
	ipInfo.IP = ipStr
	err = nil
	return
}

func (this *Ip2Region) BtreeSearch(ipStr string) (ipInfo IpInfo, err error) {
	ipInfo = IpInfo{IP: ipStr}
	ip, err := ip2long(ipStr)

	if this.headerLen == 0 {
		this.dbFileHandler.Seek(8, 0)

		buffer := make([]byte, TOTAL_HEADER_LENGTH)
		this.dbFileHandler.Read(buffer)
		var idx int64
		for i := 0; i < TOTAL_HEADER_LENGTH; i += 8 {
			startIp := getLong(buffer, int64(i))
			dataPar := getLong(buffer, int64(i+4))
			if dataPar == 0 {
				break
			}

			this.headerSip = append(this.headerSip, startIp)
			this.headerPtr = append(this.headerPtr, dataPar)
			idx++
		}

		this.headerLen = idx
	}

	var l, sptr, eptr int64
	h := this.headerLen

	for l <= h {
		m := int64(l+h) >> 1
		if m < this.headerLen {
			if ip == this.headerSip[m] {
				if m > 0 {
					sptr = this.headerPtr[m-1]
					eptr = this.headerPtr[m]
				} else {
					sptr = this.headerPtr[m]
					eptr = this.headerPtr[m+1]
				}
				break
			}
			if ip < this.headerSip[m] {
				if m == 0 {
					sptr = this.headerPtr[m]
					eptr = this.headerPtr[m+1]
					break
				} else if ip > this.headerSip[m-1] {
					sptr = this.headerPtr[m-1]
					eptr = this.headerPtr[m]
					break
				}
				h = m - 1
			} else {
				if m == this.headerLen-1 {
					sptr = this.headerPtr[m-1]
					eptr = this.headerPtr[m]
					break
				} else if ip <= this.headerSip[m+1] {
					sptr = this.headerPtr[m]
					eptr = this.headerPtr[m+1]
					break
				}
				l = m + 1
			}
		}

	}

	if sptr == 0 {
		err = errors.New("not found")
		return
	}

	blockLen := eptr - sptr
	this.dbFileHandler.Seek(sptr, 0)
	index := make([]byte, blockLen+INDEX_BLOCK_LENGTH)
	this.dbFileHandler.Read(index)
	var dataptr int64
	h = blockLen / INDEX_BLOCK_LENGTH
	l = 0

	for l <= h {
		m := int64(l+h) >> 1
		p := m * INDEX_BLOCK_LENGTH
		sip := getLong(index, p)
		if ip < sip {
			h = m - 1
		} else {
			eip := getLong(index, p+4)
			if ip > eip {
				l = m + 1
			} else {
				dataptr = getLong(index, p+8)
				break
			}
		}
	}

	if dataptr == 0 {
		err = errors.New("not found")
		return
	}

	dataLen := (dataptr >> 24) & 0xFF
	dataPtr := dataptr & 0x00FFFFFF

	this.dbFileHandler.Seek(dataPtr, 0)
	data := make([]byte, dataLen)
	this.dbFileHandler.Read(data)
	ipInfo = getIpInfo(getLong(data, 0), data[4:])
	ipInfo.IP = ipStr
	return
}

func getLong(b []byte, offset int64) int64 {

	val := int64(b[offset]) |
		int64(b[offset+1])<<8 |
		int64(b[offset+2])<<16 |
		int64(b[offset+3])<<24
	return val
}

func ip2long(IpStr string) (int64, error) {
	bits := strings.Split(IpStr, ".")
	if len(bits) != 4 {
		return 0, errors.New("ip format error")
	}

	var sum int64
	for i, n := range bits {
		bit, _ := strconv.ParseInt(n, 10, 64)
		sum += bit << uint(24-8*i)
	}

	return sum, nil
}
