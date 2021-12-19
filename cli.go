package meiwobuxing

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/LyricTian/queue"
	. "github.com/r3inbowari/zlog"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var tag = "CLI"

var CliPath = "cli.exe"

var (
	ErrUnknown       = errors.New("unknown error")
	ErrEmptyUserInfo = errors.New("empty user info")
	ErrErrorUserInfo = errors.New("error user info")
	ErrOpenCLI       = errors.New("open failed")
	ErrLogout        = errors.New("logout failed")
	ErrNonExist      = errors.New("non exist dir")
	ErrCopyAndroid   = errors.New("安卓文件转移失败")
	ErrCopyIOS       = errors.New("IOS文件转移失败")
	ErrCopyProject   = errors.New("项目文件转移失败")
)

var CliQueue *queue.Queue // 打包信号
var PackBuildMap sync.Map // 构建信息地图
const CliThread = 1       // 脚手架数量

const (
	StateHangUp  = 0
	StateRunning = 1
	StateOk      = 2
	StateFailed  = -1
)

// GetCliVersion 脚手架版本查询
func GetCliVersion() string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cmd := exec.CommandContext(ctx, CliPath, "ver")
	output, err := cmd.Output()
	if err != nil || ctx.Err() == context.DeadlineExceeded {
		return ""
	}
	return string(output)
}

func CliUserLogin(username, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if len(username) == 0 || len(password) == 0 {
		return ErrEmptyUserInfo
	}
	cmd := exec.CommandContext(ctx, CliPath, "user", "login", "--username", username, "--password", password)
	output, err := cmd.Output()
	if ctx.Err() == context.DeadlineExceeded {
		return ctx.Err()
	}
	if err != nil {
		return err
	}
	result := string(output)
	if strings.Contains(result, "0:user login:OK") {
		return nil
	} else if strings.Contains(result, "不匹配") {
		return ErrErrorUserInfo
	}
	return ErrUnknown
}

func VerifyUser(username string) error {
	un, err := CliUserInfo()
	if err == context.DeadlineExceeded {
		return err
	}
	if un == username || un == username+"\r" {
		return nil
	}
	return errors.New("verify user failed: current" + un)
}

func CliUserInfo() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cmd := exec.CommandContext(ctx, CliPath, "user", "info")
	output, err := cmd.Output()
	if ctx.Err() == context.DeadlineExceeded {
		return "", ctx.Err()
	}
	if err != nil {
		return "", err
	}
	result := string(output)
	arr := strings.Split(result, "\n")

	if !strings.Contains(arr[1], "0:user info:OK") {
		return "", ErrUnknown
	}
	return arr[0], nil
}

func CliUserLogout() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cmd := exec.CommandContext(ctx, CliPath, "user", "logout")
	output, err := cmd.Output()
	println(string(output))
	if ctx.Err() == context.DeadlineExceeded {
		return ctx.Err()
	}
	if err != nil {
		return err
	}
	result := string(output)
	if strings.Contains(result, "0:user logout:OK") {
		return nil
	} else {
		return ErrLogout
	}
}

func CliOpenProject(projectName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cmd := exec.CommandContext(ctx, CliPath, "project", "open", "--path", Up.RunPath+"/"+projectName)
	_, err := cmd.Output()
	if ctx.Err() == context.DeadlineExceeded {
		return ctx.Err()
	}
	if err != nil {
		return err
	}
	return err
}

func CliCloseProject(projectName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cmd := exec.CommandContext(ctx, CliPath, "project", "close", "--path", Up.RunPath+"/"+projectName)
	_, err := cmd.Output()
	if ctx.Err() == context.DeadlineExceeded {
		return ctx.Err()
	}
	if err != nil {
		return err
	}
	return err
}

type prefixSuffixSaver struct {
	N         int // max size of prefix or suffix
	prefix    []byte
	suffix    []byte // ring buffer once len(suffix) == N
	suffixOff int    // offset to write into suffix
	skipped   int64

	// TODO(bradfitz): we could keep one large []byte and use part of it for
	// the prefix, reserve space for the '... Omitting N bytes ...' message,
	// then the ring buffer suffix, and just rearrange the ring buffer
	// suffix when Bytes() is called, but it doesn't seem worth it for
	// now just for error messages. It's only ~64KB anyway.
}

// fill appends up to len(p) bytes of p to *dst, such that *dst does not
// grow larger than w.N. It returns the un-appended suffix of p.
func (w *prefixSuffixSaver) fill(dst *[]byte, p []byte) (pRemain []byte) {
	if remain := w.N - len(*dst); remain > 0 {
		add := minInt(len(p), remain)
		*dst = append(*dst, p[:add]...)
		p = p[add:]
	}
	return p
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (w *prefixSuffixSaver) Write(p []byte) (n int, err error) {
	lenp := len(p)
	p = w.fill(&w.prefix, p)

	// Only keep the last w.N bytes of suffix data.
	if overage := len(p) - w.N; overage > 0 {
		p = p[overage:]
		w.skipped += int64(overage)
	}
	p = w.fill(&w.suffix, p)

	// w.suffix is full now if p is non-empty. Overwrite it in a circle.
	for len(p) > 0 { // 0, 1, or 2 iterations.
		n := copy(w.suffix[w.suffixOff:], p)
		p = p[n:]
		w.skipped += int64(n)
		w.suffixOff += n
		if w.suffixOff == w.N {
			w.suffixOff = 0
		}
	}
	return lenp, nil
}

func (w *prefixSuffixSaver) Bytes() []byte {
	if w.suffix == nil {
		return w.prefix
	}
	if w.skipped == 0 {
		return append(w.prefix, w.suffix...)
	}
	var buf bytes.Buffer
	buf.Grow(len(w.prefix) + len(w.suffix) + 50)
	buf.Write(w.prefix)
	buf.WriteString("\n... omitting ")
	buf.WriteString(strconv.FormatInt(w.skipped, 10))
	buf.WriteString(" bytes ...\n")
	buf.Write(w.suffix[w.suffixOff:])
	buf.Write(w.suffix[:w.suffixOff])
	return buf.Bytes()
}

func (c *PackRequest) execCommand(commandName string, params ...string) (string, error) {
	// 打包的时候允许10分钟超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	cmd := exec.CommandContext(ctx, commandName, params...)
	// 显示运行的命令
	fmt.Println(cmd.Args)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "null", err
	}
	// 运行
	err = cmd.Start()
	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(stdout)
	defer stdout.Close()
	var retStr string
	for {
		line, err2 := reader.ReadBytes('\n')
		if err2 != nil || io.EOF == err2 {
			Log.WithTag(tag).WithFields(logrus.Fields{"err": err2.Error()}).Info("已退出控制台")
			break
		}
		lineStr := ConvertByte2StringInWindows(line, GB18030)
		retStr += lineStr
		// TODO bug builtin println 可能导致死锁
		fmt.Print(lineStr)
		// 实时添加到log中并刷新缓存
		c.Log = retStr
		c.SetBuildState(StateRunning)
		if strings.Contains(retStr, "错误") {
			err = errors.New("打包发生错误，请查看日志")
		}
	}
	cmd.Wait()
	if ctx.Err() == context.DeadlineExceeded {
		return retStr, ctx.Err()
	}
	return retStr, err
}

func (c *PackRequest) CliPack() error {
	output, err := c.execCommand(CliPath, "pack", "--config", GetPath(c.TraceID))
	if err != nil {
		return err
	}
	if Count(output) {
		return errors.New("the number of packaging has been exhausted")
	}
	c.DownloadLink, err = GetHyperLinks(output)
	// c.Log = output
	return err
}

func Count(s string) bool {
	return strings.Contains(s, "明天再来")
}

// CliOpen TODO cli may be clash
func CliOpen() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	cmd := exec.CommandContext(ctx, CliPath, "open")
	_, err := cmd.Output()
	if ctx.Err() == context.DeadlineExceeded {
		return ctx.Err()
	}
	return err
}

func (c *PackRequest) SetBuildState(t int) {
	c.State = t
	PackBuildMap.Store(c.TraceID, c)
}

func GetBuildState(traceID string) *PackRequest {
	value, ok := PackBuildMap.Load(traceID)
	if !ok {
		return nil
	}
	return value.(*PackRequest)
}

func InitCli() {
	Log.WithTag(tag).Info("setting HBuilderX environment variables")
	if runtime.GOOS == "darwin" {
		CliPath = "/Applications/HBuilderX.app/Contents/MacOS/cli"
	}

	Log.WithTag(tag).Info("starting HBuilderX")
	err := CliOpen()
	if err != nil {
		Log.WithTag(tag).Error("please make sure HBuilderX is installed or the execution path is correct")
		Log.WithTag(tag).WithFields(logrus.Fields{"err": err.Error()}).Info("exit")
		time.Sleep(time.Minute)
		os.Exit(1004)
	}
	_ = KillProcess(HBuilderXTargetName)

	// 善后处理
	AddExitFunc("HbuilderX", func() {
		err := KillProcess(HBuilderXTargetName)
		if err != nil {
			Log.WithTag(tag).Warn("panic during kill HBuilderX")
			return
		}
	})
	// 初始化并启动任务队列
	CliQueue = queue.NewQueue(100000, CliThread)
	CliQueue.Run()
}

var cnt = 0

func (c *PackRequest) PackEnqueue() string {
	c.TraceID = CreateUUID()
	c.SetBuildState(StateHangUp)
	cnt++
	// 打包请求入队
	job := queue.NewJob("hello", c.executePack)
	CliQueue.Push(job)
	Log.WithTag(tag).WithFields(logrus.Fields{"user": c.Username, "traceID": c.TraceID, "order": cnt}).Info("task enqueue")
	return c.TraceID
}

func (c *PackRequest) savePackConfig() error {
	return SaveJsonToRes(c.TraceID, c.PackConfig)
}

func (c *PackRequest) executePack(i interface{}) {
	// 设置为运行状态
	c.SetBuildState(StateRunning)
	// 开始
	Log.WithTag(tag).WithFields(logrus.Fields{"traceID": c.TraceID}).Info("running task")

	var err error
	// 错误处理与善后
	defer func() {
		// 附加处理原因
		c.SetBuildState(StateOk)
		if err != nil {
			time.Sleep(time.Second)
			Log.WithTag(tag).WithFields(logrus.Fields{"traceID": c.TraceID, "err": err.Error()}).Error("task done with error")
			c.Reason = err.Error()
			// 失败的话再次关闭
			_ = CliCloseProject(c.Package)
			c.SetBuildState(StateFailed)
		} else {
			Log.WithTag(tag).WithFields(logrus.Fields{"traceID": c.TraceID}).Info("task done")
		}
		// 过程结束调用回调接口
		Log.WithTag(tag).WithFields(logrus.Fields{"traceID": c.TraceID}).Info("all done...bye (●’◡’●)ﾉ notify -> " + c.Callback)
		Get(c.Callback + "?trace=" + c.TraceID)
		// 关闭 HBX
		time.Sleep(time.Second * 1)
		_ = KillProcess(HBuilderXTargetName)
	}()

	t := NewTask()

	t.AddProcess("打开 HBuilderX", CliOpen, 2)
	//t.AddProcess("退出登陆", CliUserLogout, 10)
	//t.AddProcess("是否退出登录", func() error {
	//	return VerifyUser("")
	//}, 20)
	t.AddProcess("用户登录", func() error {
		return CliUserLogin(c.Username, c.Password)
	}, 5)
	t.AddProcess("用户校验", func() error {
		return VerifyUser(c.Username)
	}, 2)
	t.AddProcess("Android 证书替换", func() error {
		// keystore
		if VerifyUUID(c.Certs.Android) {
			c.PackConfig.Android.CertFile = GetPath(c.Certs.Android)
		}
		return nil
	})
	t.AddProcess("IOS 证书替换", func() error {
		// p12/mobileprovision
		if VerifyUUID(c.Certs.IOS) && VerifyUUID(c.Certs.IOSEx) {
			c.PackConfig.IOS.CertFile = GetPath(c.Certs.IOS)
			c.PackConfig.IOS.Profile = GetPath(c.Certs.IOSEx)
		}
		return nil
	})
	t.AddProcess("保存打包配置文件", c.savePackConfig)
	t.AddProcess("查询目标项目是否存在", func() error {
		if !Exists(c.Package + "/") {
			return ErrNonExist
		}
		return nil
	})
	t.AddProcess("打开项目", func() error {
		return CliOpenProject(c.Package)
	}, 1)
	t.AddProcess("Android 安全图拷贝", func() error {
		// 拷贝安卓安全图
		androidPath := Up.RunPath + "/" + c.Package + "/nativeplugins/Html5app-Baichuan/android/res/drawable/yw_1222.jpg"
		if CopyFormRes(c.Safe.Android, androidPath) != nil {
			return ErrCopyAndroid
		}
		return nil
	})
	t.AddProcess("IOS 安全图拷贝", func() error {
		// 拷贝IOS安全图
		iosPath := Up.RunPath + "/" + c.Package + "/nativeplugins/Html5app-Baichuan/ios/yw_1222.jpg"
		if CopyFormRes(c.Safe.IOS, iosPath) != nil {
			return ErrCopyIOS
		}
		return nil
	})
	t.AddProcess("Manifest 转换", func() error {
		manifest, err1 := OpenManifest(c.Manifest)
		if err1 != nil {
			return err1
		}
		// 图标路径转换更新
		err1 = manifest.ConvertPath(c.Package)
		if err1 != nil {
			return err1
		}
		// 拷贝到项目
		manifestPath := Up.RunPath + "/" + c.Package + "/manifest.json"
		if CopyFormRes(c.Manifest, manifestPath) != nil {
			err1 = ErrCopyProject
			return err1
		}
		return nil
	})

	t.AddProcess("App 环境变量导入", func() error {
		if !VerifyUUID(c.Export) {
			Log.WithTag(tag).Warn("未导入环境变量 config.js")
			return nil
		}
		exportPath := Up.RunPath + "/" + c.Package + "/utils/config.js"
		if CopyFormRes(c.Export, exportPath) != nil {
			return ErrCopyProject
		}
		return nil
	})

	t.AddProcess("启动云构建", func() error {
		// 11. 调用云构建
		return c.CliPack()
	}, 2)

	t.AddProcess("清理现场", func() error {
		return CliCloseProject(c.Package)
	}, 2)

	err = t.Start()
}
