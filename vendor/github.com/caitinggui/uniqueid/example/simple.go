package main

import (
	"fmt"
	"hash/fnv"
	"os"
	"time"

	"github.com/caitinggui/uniqueid"
)

func main() {
	// 此处也可从MySQL或者zookeeper等获取到当前的worker id

	hostname, _ := os.Hostname() // docker启动后的hostname也是不同的
	hasher := fnv.New32a()
	hasher.Write([]byte(hostname))
	workerId := uint16(hasher.Sum32() & (1<<10 - 1)) // 这里workerId是11位的
	var reserveId uint8 = 0

	sf := uniqueid.NewUniqueId(workerId, reserveId)
	for i := 0; i < 20; i++ {
		uid, err := sf.NextId()
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
		fmt.Println("uid: ", uid, uniqueid.Prase(uid))
	}
}
