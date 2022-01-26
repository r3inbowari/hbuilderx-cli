package meiwobuxing

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

var HBuilderXTargetName = "HbuilderX"

func runInWindows(cmd string) (string, error) {
	result, err := exec.Command("cmd", "/c", cmd).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(result)), err
}

func RunCommand(cmd string) (string, error) {
	if runtime.GOOS == "windows" {
		return runInWindows(cmd)
	} else {
		return runInLinux(cmd)
	}
}

func runInLinux(cmd string) (string, error) {
	fmt.Println("Running Linux cmd:" + cmd)
	result, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(result)), err
}

// CheckProRunning 根据进程名判断进程是否运行
func CheckProRunning(serverName string) (bool, error) {
	a := `ps ux | awk '/` + serverName + `/ && !/awk/ {print $2}'`
	pid, err := RunCommand(a)
	if err != nil {
		return false, err
	}
	return pid != "", nil
}

// GetPid 根据进程名称获取进程ID
func GetPids(serverName string) ([]string, error) {
	a := `ps ux | awk '/` + serverName + `/ && !/awk/ {print $2}'`
	pid, err := RunCommand(a)
	pids := strings.Split(pid, "\n")
	if pids[0] == "" {
		return nil, nil
	}
	return pids, err
}

func KillProcess(target string) error {
	pids, err := GetPids(target)
	if err != nil {
		return err
	}
	// kill it
	fmt.Printf("[PD] not supported method, len: %d\n", len(pids))
	return nil
}
