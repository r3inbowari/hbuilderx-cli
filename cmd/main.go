package main

import (
	"fmt"
	"github.com/hpcloud/tail"
	"gopkg.in/alecthomas/kingpin.v2"
	"meiwobuxing"
	"os"
	"os/user"
	"runtime"
	"time"
)

var (
	md       = kingpin.Flag("md5", "create md5 using username and password").Short('m').Bool()
	username = kingpin.Arg("username", "first arg: used to create md5").Default("caicai").String()
	password = kingpin.Arg("password", "second arg: used to create md5").Default("159463").String()
)

var (
	gitHash        string
	buildTime      string
	goVersion      string
	releaseVersion string
	major          string
	minor          string
	patch          string
)

var mode = "DEV"

func main() {
	kingpin.Version(fmt.Sprintf("%s git-%s build on %s", releaseVersion, gitHash, buildTime))
	kingpin.Parse()
	if *md {
		fmt.Printf("[MD5] username: %s\n", *username)
		fmt.Printf("[MD5] password: %s\n", *password)
		fmt.Printf("[MD5] computed: %s\n", meiwobuxing.CalcJwtMD5(*username, *password))
		return
	}
	// 退出信号拦截
	meiwobuxing.InitSignalExit(func(signal os.Signal) {})
	// 开发参数注入
	dev()
	// 更新插件
	meiwobuxing.InitUpdate(buildTime, mode, releaseVersion, gitHash, major, minor, patch, "meiwobuxing", nil)
	// 文件系统初始化
	meiwobuxing.InitFileSystem(meiwobuxing.Up.RunPath, 100)
	// 配置插件
	meiwobuxing.InitConfig()
	// 日志
	meiwobuxing.InitLogger(meiwobuxing.Up.BuildMode, OpenTracker)
	// 更新服务与权限
	meiwobuxing.SoftwareUpdate(false)
	// 脚手架初始化
	meiwobuxing.InitCli()
	// 服务启动
	meiwobuxing.CLIApplication()
}

func dev() {
	if mode == "DEV" {
		buildTime = "Thu Oct 01 00:00:00 1970 +0800"
		gitHash = "cb0dc838e04e841f193f383e06e9d25a534c5809"
		goVersion = runtime.Version()
		releaseVersion = "ver[DEV]"
	}
}

func OpenTracker() {
	fileName := ""
	u, err := user.Current()
	if meiwobuxing.GetConfig(false).LogPath == "" {
		fileName = u.HomeDir + "/AppData/Roaming/HBuilder X/.log"
	} else {
		fileName = meiwobuxing.GetConfig(false).LogPath
	}
	if !meiwobuxing.Exists(fileName) {
		meiwobuxing.Log.Warn("[Track] not found log")
		return
	}
	config := tail.Config{
		ReOpen:    true,                                 // 重新打开
		Follow:    true,                                 // 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 0}, // Seek
		MustExist: false,                                // 必须存在
		Poll:      true,
		Logger:    meiwobuxing.Log,
	}
	tails, err := tail.TailFile(fileName, config)
	if err != nil {
		meiwobuxing.Log.Error("[Tracker] tail file failed, err:", err)
		return
	}
	var (
		line *tail.Line
		ok   bool
	)
	for {
		line, ok = <-tails.Lines
		if !ok {
			meiwobuxing.Log.Errorf("[Tracker] tail file close reopen, filename:%s\n", tails.Filename)
			time.Sleep(time.Second)
			continue
		}
		fmt.Println("[T] [Tracker] line:", line.Text)
	}
}
