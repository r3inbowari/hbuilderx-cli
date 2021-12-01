package test

import (
	"fmt"
	"github.com/LyricTian/queue"
	"log"
	"meiwobuxing"
	"os/exec"
	"testing"
	"time"
)

func TestGetVersion(t *testing.T) {
	meiwobuxing.CliPath = "/Applications/HBuilderX.app/Contents/MacOS/cli"
	println(meiwobuxing.GetCliVersion())
	time.Sleep(time.Second)
}

func TestCliUserLogin(t *testing.T) {
	meiwobuxing.CliPath = "/Applications/HBuilderX.app/Contents/MacOS/cli"
	err := meiwobuxing.CliUserLogin("212156620@qq.com", "osjava.cn")
	println(err)
	time.Sleep(time.Second)
	err1 := meiwobuxing.CliUserLogin("212156620@qq.com", "osjava.c")
	println(err1)
	time.Sleep(time.Second)
	println(meiwobuxing.VerifyUser("212156620@qq.com"))
	time.Sleep(time.Second)
}

func TestCliUserInfo(t *testing.T) {
	meiwobuxing.CliPath = "/Applications/HBuilderX.app/Contents/MacOS/cli"
	println(meiwobuxing.CliUserInfo())
	time.Sleep(time.Second)
}

func TestCliUserLogout(t *testing.T) {
	meiwobuxing.CliPath = "/Applications/HBuilderX.app/Contents/MacOS/cli"
	println(meiwobuxing.CliUserLogout())
	time.Sleep(time.Second)
}

func TestOpen(t *testing.T) {
	meiwobuxing.CliPath = "/Applications/HBuilderX.app/Contents/MacOS/cli"
	println(meiwobuxing.CliOpen())
	time.Sleep(time.Second)
}

func TestPack(t *testing.T) {
	//meiwobuxing.CliPath = "/Applications/HBuilderX.app/Contents/MacOS/cli"
	//time.Sleep(time.Second * 1)
	//var cli meiwobuxing.PackRequest
	//cli.TraceID = "configure"
	//cli.CliPack()
	//time.Sleep(time.Second * 1)
}

func TestStateSG(t *testing.T) {
	var cli meiwobuxing.PackRequest
	cli.TraceID = "HBuilderX"
	cli.Reason = "goto"

	cli.SetBuildState(meiwobuxing.StateOk)
	ps := meiwobuxing.GetBuildState("HBuilderX")
	println(ps.Reason)

	ps = meiwobuxing.GetBuildState("HBuilder")
	println(ps == nil)
}

func TestCMD(t *testing.T) {
	RunCommand("ping", "192.168.5.1")
}

func RunCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	log.Println(cmd.String())
	// 命令的错误输出和标准输出都连接到同一个管道
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	if err != nil {
		return err
	}
	if err = cmd.Start(); err != nil {
		return err
	}
	// 从管道中实时获取输出并打印到终端
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		fmt.Print(string(tmp))
		if err != nil {
			break
		}
	}
	if err = cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func TestCli(t *testing.T) {
	q := queue.NewQueue(1, 1)
	q.Run()

	defer q.Terminate()

	job := queue.NewJob("hello", func(v interface{}) {
		fmt.Printf("%s,world \n", v)
		time.Sleep(time.Second*5)
	})
	q.Push(job)

	for i := 0; i < 100; i++ {
		fc(q, i)
	}

}

func fc(q *queue.Queue, i int) {
	job1 := queue.NewJob("hello", func(v interface{}) {
		//time.Sleep(time.Second*3)
		print(i)
		fmt.Printf("%s,world1 \n", v)
		//time.Sleep(time.Second)

	})
	q.Push(job1)
}

func TestCliBuild(t *testing.T) {
	//go func() {
	//	c := exec.Command("/Applications/HBuilderX.app/Contents/MacOS/cli", "pack", "--config", "/Users/r3inb/Downloads/meiwobuxing/res/52acad9c87f09fda8e1f004a4c47b818/81c6a952-d955-441b-9232-2bcc36db3c95")
	//	//o, err := c.Output()
	//	//if err != nil {
	//	//	println("er")
	//	//}
	//	//println(string(o))
	//	c.Run()
	//
	//	println("aaa")
	//}()
	//
	//time.Sleep(time.Second*100)
}