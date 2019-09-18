package process

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"
)

const (
	SIG_SYS_READ = syscall.Signal(0)
)

func IsRunningUsingFile(pidFile string) (bool, error) {
	_, err := os.Stat(pidFile)
	if os.IsNotExist(err) {
		// pid file is not exist, process is finished
		return false, nil
	}

	contents, err := ioutil.ReadFile(pidFile)
	if err != nil {
		return false, err
	}

	pid, err := strconv.Atoi(string(contents))
	if err != nil {
		return false, err
	}

	return IsRunning(pid)
}

func IsRunning(pid int) (bool, error) {
	switch runtime.GOOS {
	case "linux":
		return isRunningOnLinux(pid)
	default:
		return false, errors.New(fmt.Sprintf("Unsupported operating system: %s", runtime.GOOS))
	}

}

func WritePIDFile(pidFile string) error {
	dirPath := filepath.Dir(pidFile)
	if dirPath != "" {
		os.MkdirAll(dirPath, os.FileMode(0700))
	}

	return ioutil.WriteFile(pidFile, []byte(fmt.Sprintf("%d", os.Getpid())), os.FileMode(0600))
}

func isRunningOnLinux(pid int) (bool, error) {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false, err
	}

	err = process.Signal(SIG_SYS_READ)
	if err == nil {
		// process is running
		return true, nil
	}

	switch err.(type) {
	case *os.SyscallError:
		// signal call not allowed
		return true, nil
	}

	// process is finished
	return false, nil
}
