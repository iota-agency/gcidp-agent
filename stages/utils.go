package stages

import (
	"errors"
	"os/exec"
	"strings"
)

func runCmd(cmd *exec.Cmd) error {
	var outbuf, errbuf strings.Builder
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	if err := cmd.Run(); err != nil {
		return errors.New(errbuf.String())
	}
	return nil
}
