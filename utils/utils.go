package utils

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
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

func BuildPlugin(dst string, src string) error {
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", dst, src)
	return RunCmd(cmd)
}

type LogLine struct {
	Stream string `json:"stream"`
}

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

func Print(rd io.Reader) error {
	var lastLine string

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		lastLine = scanner.Text()
		log := &LogLine{}
		if err := json.Unmarshal(scanner.Bytes(), log); err != nil {
			return err
		}
	}

	errLine := &ErrorLine{}
	if err := json.Unmarshal([]byte(lastLine), errLine); err != nil {
		return err
	}
	if errLine.Error != "" {
		return errors.New(errLine.Error)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
