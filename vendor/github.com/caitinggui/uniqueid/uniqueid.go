package uniqueid

import (
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"
)

const (
	nano              = 1000 * 1000 // 纳秒转为毫秒
	BitLenTime        = 41          // 41位存时间戳，大约可用69年
	BitLenWorker      = 10          // 11位存机器id
	BitLenReserve     = 2           // 2位预留，可用于业务编码
	BitLenAbnormality = 2           // 2位用于时间回退等异常情况，运行时可恢复
	BitLenSequence    = 8           // 8位用户1毫秒内的递增值，意味着系统1s内并发在(2^8-1)*1000=255000
	BitLenTotal       = BitLenTime + BitLenWorker + BitLenReserve + BitLenAbnormality + BitLenSequence
)

var (
	rander       = rand.New(rand.NewSource(time.Now().UnixNano()))
	sequenceMask = uint64(1<<BitLenSequence - 1)
	StartTime    = time.Date(2019, 1, 25, 0, 0, 0, 0, time.UTC).UnixNano() / nano // 起始时间，更改可能会导致id重复

)

type UniqueId struct {
	mutex     *sync.Mutex
	startTime int64
	lastTime  int64
	errorTime int64

	workerId      uint16
	reserveId     uint8
	abnormalityId uint8
	sequenceId    uint64
}

// 生成下一个unique id，如果该毫秒内的id数超过sequence的最大值，就报错
func (self *UniqueId) NextId() (uid uint64, err error) {
	// 加锁，保证sequence的串行
	self.mutex.Lock()
	current := currentMillisecond()
	if current > self.lastTime {
		// 时间戳到了下一个毫秒后, sequence 初始值用一个随机数，避免最后一位都是同一样
		self.sequenceId = uint64(rander.Intn(10)) // 返回10以内的随机整数, 不会太耗性能

	} else if current == self.lastTime {
		// 时间戳在同一个毫秒内,sequence加1同时判断是否溢出
		self.sequenceId = (self.sequenceId + 1) & sequenceMask
		// 溢出了
		if self.sequenceId == 0 {
			log.Println("sequence overflow")
			time.Sleep(time.Millisecond)   // 休息1毫秒
			current = currentMillisecond() //启用新的时间戳, 此时sequence正好为0

		}

	} else {
		// 当前时间突然小于上一次的时间，说明服务器时间有倒退（可能是同步ntp引起）
		// 这里碰到错误就把异常码加1，同时要重置lastTime
		// 如果服务器时间被重置了？如果重置了一次，然后还没到正确时间，又被重置了呢？ 这种方法可以解决
		self.abnormalityId = (self.abnormalityId + 1) & (1<<BitLenAbnormality - 1)
		if self.abnormalityId == 0 {
			self.mutex.Unlock()                // 比用defer快一些
			return 0, errors.New("time error") // 错误码容量超出

		}
		if self.errorTime < self.lastTime {
			self.errorTime = self.lastTime // errorTime总是取最大值

		}
		self.lastTime = current

	}

	// 时间要统一更新一下
	self.lastTime = current
	// 时间已追平，恢复错误码
	if self.errorTime != 0 && self.errorTime < self.lastTime {
		self.errorTime = 0
		self.abnormalityId = 0

	}
	uid = self.toId()
	self.mutex.Unlock() // 比用defer快一些
	return

}

// 不考虑时间戳超过最大值，毕竟是60多年之后的事情
func (self *UniqueId) toId() uint64 {
	return uint64(self.lastTime)<<(BitLenTotal-BitLenTime) |
		uint64(self.workerId)<<(BitLenTotal-BitLenTime-BitLenWorker) |
		uint64(self.reserveId)<<(BitLenSequence+BitLenAbnormality) |
		uint64(self.abnormalityId)<<BitLenSequence |
		uint64(self.sequenceId)

}

// 解析id的含义
func Prase(uid uint64) map[string]uint64 {
	workerMask := uint64(1<<(BitLenTotal-BitLenTime) - 1)
	reserveMask := uint64(1<<(BitLenSequence+BitLenAbnormality+BitLenReserve) - 1)
	abnormalityMask := uint64(1<<(BitLenSequence+BitLenAbnormality) - 1)

	lastTime := uid >> (BitLenTotal - BitLenTime)
	workerId := (uid & workerMask) >> (BitLenTotal - BitLenTime - BitLenWorker)
	reserveId := (uid & reserveMask) >> (BitLenSequence + BitLenAbnormality)
	abnormalityId := (uid & abnormalityMask) >> BitLenSequence
	sequenceId := uid & uint64(sequenceMask)
	return map[string]uint64{
		"time":          lastTime,
		"workerId":      workerId,
		"reserveId":     reserveId,
		"abnormalityId": abnormalityId,
		"sequenceId":    sequenceId,
	}

}

// 返回id生成器的实例
// 配置 WorkerId 机器id; ReserveId 预留的值，可以是业务编码
func NewUniqueId(WorkerId uint16, ReserveId uint8) *UniqueId {
	var (
		sf           UniqueId
		maxWorkerId  uint16 = 1<<BitLenWorker - 1
		maxReserveId uint8  = 1<<BitLenReserve - 1
	)
	// 参数不合法
	if WorkerId > maxWorkerId || ReserveId > maxReserveId {
		panic("invalid parameter")

	}
	sf.workerId = WorkerId
	sf.reserveId = ReserveId
	sf.lastTime = currentMillisecond()
	sf.mutex = new(sync.Mutex)
	return &sf

}

// 返回当前的毫秒时间戳
func currentMillisecond() int64 {
	return time.Now().UnixNano() / nano

}
