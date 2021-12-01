package meiwobuxing

import (
	"bytes"
	"errors"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
)

var HBuilderXTargetName = "HbuilderX"

const TH32CS_SNAPPROCESS = 0x00000002

type WindowsProcess struct {
	ProcessID       int
	ParentProcessID int
	Exe             string
}

func newWindowsProcess(e *syscall.ProcessEntry32) WindowsProcess {
	end := 0
	for {
		if e.ExeFile[end] == 0 {
			break
		}
		end++
	}
	return WindowsProcess{
		ProcessID:       int(e.ProcessID),
		ParentProcessID: int(e.ParentProcessID),
		Exe:             syscall.UTF16ToString(e.ExeFile[:end]),
	}
}

func Processes() ([]WindowsProcess, error) {
	handle, err := syscall.CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer syscall.CloseHandle(handle)

	var entry syscall.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))
	// get the first process
	err = syscall.Process32First(handle, &entry)
	if err != nil {
		return nil, err
	}

	results := make([]WindowsProcess, 0, 50)
	for {
		results = append(results, newWindowsProcess(&entry))

		err = syscall.Process32Next(handle, &entry)
		if err != nil {
			// windows sends ERROR_NO_MORE_FILES on last process
			if err == syscall.ERROR_NO_MORE_FILES {
				return results, nil
			}
			return nil, err
		}
	}
}

func FindProcessByName(processes []WindowsProcess, name string) *WindowsProcess {
	for _, p := range processes {
		if bytes.Contains([]byte(strings.ToUpper(p.Exe)), []byte(strings.ToUpper(name))) {
			return &p
		}
	}
	return nil
}

func KillProcess(target string) error {
	if runtime.GOOS == "windows" {
		processes, err := Processes()
		if err != nil {
			return err
		}
		p := FindProcessByName(processes, target)
		if p != nil {
			c := exec.Command("taskkill.exe", "/f", "/im", p.Exe)
			return c.Start()
		}
		return errors.New("not found process: " + target)
	} else {
		return errors.New("not supported methods")
	}
}
