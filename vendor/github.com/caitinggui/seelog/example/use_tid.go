package main

import (
	"runtime"
	"time"

	logger "github.com/caitinggui/seelog"
)

var logConfig = `
<seelog minlevel="trace">
  <outputs formatid="common">

    <!--用来记录gin的log-->
    <filter levels="trace">
      <console formatid="access" />
      <rollingfile formatid="access" type="size" filename="logs/access_log.log" maxsize="102400000" maxrolls="7"/>
    </filter>

    <filter levels="info,debug,warn,error,critical">
      <console />
      <!--在golang中，2006代表年份，01代表月份，02代表日，因为golang是2016-01-02诞生的-->
      <rollingfile type="date" filename="logs/[%ProgramName].log" datepattern="2006.01.02" maxrolls="7"/>
    </filter>
    <filter levels="error,critical">
      <!--maxsize单位是字节，102400大概是100k大小 过期的错误日志会被压缩保存-->
      <rollingfile type="size" filename="logs/[%ProgramName].error" maxsize="10240" archivetype="gzip" archivepath="logs/[%ProgramName].error.bak.tar.gz" maxrolls="7"/>
    </filter>

  </outputs>
  <formats>
    <format id="common" format="%Date(2006-01-02 15:04:05.000) [%LEV][%File:%Line][%Func][%Tid] %Msg%n" />
    <format id="access" format="[%Tid] %Msg" />
  </formats>
</seelog>
`

func testLogger() {
	logger.Info("hello: ", runtime.Goid())
	logger.Error("hello: ", runtime.Goid())
}

func main() {
	log, err := logger.LoggerFromConfigAsString(logConfig)

	if err != nil {
		panic(err)
	}

	// 整个进程的logger都是用log了
	logger.ReplaceLogger(log)
	// 不要忘记flush
	defer logger.Flush()

	go logger.Info("hello info")
	go logger.Trace("hello trace")
	for i := 1; i < 1000; i++ {
		go testLogger()

	}
	time.Sleep(time.Second)
	logger.Error(logger.HideStr("hello error"))
	logger.Info(logger.HideStr("hello info"))
}
