package test

import (
	"meiwobuxing"
	"os/exec"
	"testing"
)

func TestProcName(t *testing.T) {
	//pid, err := meiwobuxing.GetPid("HBuilderX.exe")
	//if err != nil {
	//	return
	//}
	//
	//println(len(pid))
}

func TestProc(t *testing.T) {
	processes, err := meiwobuxing.Processes()
	if err != nil {
		return
	}
	d := meiwobuxing.FindProcessByName(processes, "HbuilderX")
	if d != nil {
		println(d.Exe)

		c := exec.Command("taskkill.exe", "/f", "/im", d.Exe)
		err = c.Start()
		if err != nil {
			print(err.Error())
			return
		}
	}

	err = meiwobuxing.KillProcess(meiwobuxing.HBuilderXTargetName)
	if err != nil {
		println(err.Error())
		return
	}

}
