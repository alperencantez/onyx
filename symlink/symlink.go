package symlink

import (
	"fmt"
	"os"
	"os/exec"
)

func Create(actualExecutablePath string, name string) bool {
	cmd := exec.Command("ln", "-sf", actualExecutablePath, name)
	err := cmd.Run()

	return err == nil
}

func Read(symlinkName string) bool {
	fileInfo, err := os.Lstat(symlinkName)
	if err != nil {
		fmt.Println(err)
		return false
	}

	if fileInfo.Mode()&os.ModeSymlink != 0 {
		_, err := os.Readlink(fileInfo.Name())
		return err != nil
	}

	return true
}
