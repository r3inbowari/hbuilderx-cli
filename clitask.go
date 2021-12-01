package meiwobuxing

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

type CliTask struct {
	Step           int            // 当前步骤
	Processes      []func() error // 过程函数
	ProcessesName  []string       // 过程名称
	ProcessesDelay []int          // 间断时间
}

func NewTask() *CliTask {
	return &CliTask{}
}

func (c *CliTask) AddProcess(n string, f func() error, delayTime ...int) *CliTask {
	cnt := 0
	c.ProcessesName = append(c.ProcessesName, n)
	c.Processes = append(c.Processes, f)
	if len(delayTime) > 0 {
		cnt = delayTime[0]
	}
	c.ProcessesDelay = append(c.ProcessesDelay, cnt)
	Log.WithFields(logrus.Fields{"stage": n, "delay": fmt.Sprintf("%ds", cnt)}).Info("[CLI] add stage")
	return c
}

func (c *CliTask) Start() error {
	maxRetryCount := GetConfig(false).MaxRetryCount
	// 最大超时重试 每次调用Start只允许出现 maxRetryCount 次重试
	var err error

	for i, process := range c.Processes {
		c.Step = i + 1
	retry:
		time.Sleep(time.Second * time.Duration(c.ProcessesDelay[i]))
		Log.Infof("[CLI] Stage %d: %s", c.Step, c.ProcessesName[i])
		err = process()
		// 不允许一个云构建过程超时的事务进行重试
		if err == context.DeadlineExceeded && maxRetryCount > 0 && c.ProcessesName[i] != "启动云构建" {
			Log.WithFields(logrus.Fields{"name": c.ProcessesName[i]}).Warn("[CLI] trigger the transaction retry with a stage timeout")
			// 清除错误
			err = nil
			// 重试次数减1
			maxRetryCount--
			goto retry
		}
		if err != nil {
			break
		}
	}
	return err
}
