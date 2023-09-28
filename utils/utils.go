package utils

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func RunCmd(cmd *exec.Cmd) error {
	var outbuf, errbuf strings.Builder
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	if err := cmd.Run(); err != nil {
		return errors.New(errbuf.String())
	}
	return nil
}
