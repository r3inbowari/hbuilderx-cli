package main

import (
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/r3inbowari/common"
	. "github.com/r3inbowari/zlog"
	"github.com/r3inbowari/zupdate"
	"meiwobuxing"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"time"

	"github.com/hpcloud/tail"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	md       = kingpin.Flag("md5", "create md5 using username and password").Short('m').Bool()
	username = kingpin.Arg("username", "first arg: used to create md5").Default("caicai").String()
	password = kingpin.Arg("password", "second arg: used to create md5").Default("159463").String()
)

var (
	gitHash        = "cb0dc838e04e841f193f383e06e9d25a534c5809"
	buildTime      = "Thu Oct 01 00:00:00 1970 +0800"
	goVersion      = runtime.Version()
	releaseVersion = "ver[DEV]"
	major          string
	minor          string
	patch          string
	mode           = "DEV"
)

func main() {
	kingpin.Version(fmt.Sprintf("%s git-%s build on %s", releaseVersion, gitHash, buildTime))
	kingpin.Parse()
	if *md {
		fmt.Printf("[MD5] username: %s\n", *username)
		fmt.Printf("[MD5] password: %s\n", *password)
		fmt.Printf("[MD5] computed: %s\n", meiwobuxing.CalcJwtMD5(*username, *password))
		return
	}

	meiwobuxing.Up = meiwobuxing.InitUpdate(buildTime, mode, releaseVersion, gitHash, major, minor, patch, "meiwobuxing")

	InitGlobalLogger()
	Log.SetScreen(true)

	// 权限
	perm := common.InitPermClient(common.PermOptions{
		Log:         &Log.Logger,
		CheckSource: "https://1077739472743245.cn-hangzhou.fc.aliyuncs.com/2016-08-15/proxy/perm.LATEST/perm",
		AppId:       "acd3f8c51b",
		ExpireAfter: time.Hour * 8760,
	})
	perm.Verify()

	// 更新服务
	updater := zupdate.InitUpdater(zupdate.UpdateOptions{
		EntryName:   "hbxctl",
		EntryArgs:   []string{},
		Mode:        zupdate.REL,
		Log:         &Log.Logger,
		CheckSource: "https://120.77.33.188/resources/mate.json",
	})
	ma, _ := strconv.ParseInt(major, 10, 64)
	mi, _ := strconv.ParseInt(minor, 10, 64)
	pa, _ := strconv.ParseInt(patch, 10, 64)
	updater.IncludeFile("5c74bf9c1face2dcb47bae100f2c664cdbd43404", zupdate.File{
		Name:  "hbxctl",
		Major: ma,
		Minor: mi,
		Patch: pa,
	})
	updater.CheckAndUpdate()

	if !meiwobuxing.Exists(meiwobuxing.Up.RunPath + "/certs") {
		_ = os.Mkdir(meiwobuxing.Up.RunPath+"/certs", os.ModePerm)
		_ = os.Mkdir(meiwobuxing.Up.RunPath+"/certs/android", os.ModePerm)
		_ = os.Mkdir(meiwobuxing.Up.RunPath+"/certs/ios", os.ModePerm)
	}

	time.Sleep(time.Second)
	// 退出信号拦截
	meiwobuxing.InitSignalExit(func(signal os.Signal) {})
	// 文件系统初始化
	common.InitResSystem(meiwobuxing.Up.RunPath+"/", 100)

	// 本地数据库初始化
	dbOpts := badger.DefaultOptions("./db")
	dbOpts.Logger = Log
	err := meiwobuxing.InitDB(dbOpts)
	if err != nil {
		Log.Error("[MAIN] init db failed, please feedback it to the developer.")
		time.Sleep(time.Second * 5)
		return
	}

	// 配置插件
	meiwobuxing.InitConfig()
	// 日志
	Log.SetBuildMode(common.Modes[mode]).SetRotate(meiwobuxing.GetConfig(false).RotateEnable)
	go OpenTracker()
	// 脚手架初始化
	meiwobuxing.InitCli()
	// 服务启动
	meiwobuxing.CLIApplication()
}

func OpenTracker() {
	fileName := ""
	u, err := user.Current()
	if err != nil {
		return
	}
	if meiwobuxing.GetConfig(false).LogPath == "" {
		fileName = u.HomeDir + "/AppData/Roaming/HBuilder X/.log"
	} else {
		fileName = meiwobuxing.GetConfig(false).LogPath
	}
	if !meiwobuxing.Exists(fileName) {
		Log.Warn("[Track] not found log")
		return
	}
	config := tail.Config{
		ReOpen:    true,                                 // 重新打开
		Follow:    true,                                 // 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 0}, // Seek
		MustExist: false,                                // 必须存在
		Poll:      true,
		Logger:    Log,
	}
	tails, err := tail.TailFile(fileName, config)
	if err != nil {
		Log.Error("[Tracker] tail file failed, err:", err)
		return
	}
	var (
		line *tail.Line
		ok   bool
	)
	for {
		line, ok = <-tails.Lines
		if !ok {
			Log.Errorf("[Tracker] tail file close reopen, filename:%s\n", tails.Filename)
			time.Sleep(time.Second)
			continue
		}
		fmt.Println("[T] [Tracker] line:", line.Text)
	}
}
