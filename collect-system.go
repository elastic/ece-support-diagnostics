package ecediag

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/elastic/beats/libbeat/logp"
)

func runSystemCmds(tar *Tarball) {
	l := logp.NewLogger("system_cmd")

	resc, errc := make(chan string), make(chan error)
	for _, cmd := range SystemCmd {
		go func(cmd systemCmd) {
			body, err := executeCmd(tar, cmd)
			// body, err := runCommand(cmd)
			if err != nil {
				errc <- err
				return
			}
			resc <- string(body)
		}(cmd)
	}

	for i := 0; i < len(SystemCmd); i++ {
		select {
		case res := <-resc:
			l.Info(res)
		case err := <-errc:
			l.Error(err)
		}
	}
}

func executeCmd(tar *Tarball, c systemCmd) (string, error) {
	output, err := c.run()
	if err != nil {
		err = fmt.Errorf("Command failed: `%v`, %s", c.RawCmd, err)
		return "", err
	}

	// fpath := filepath.Join(SystemTmpDir, c.Filename)
	fpath := filepath.Join(DiagName, c.Filename)
	tar.AddData(fpath, output)
	// err = writeFile(fpath, output)

	if err != nil {
		return "", err
	}
	res := fmt.Sprintf("Command completed: \"%v\"", c.RawCmd)
	return res, nil
}

func (c *systemCmd) run() ([]byte, error) {
	expand := strings.Split(c.RawCmd, " ")
	bin := expand[0]
	args := expand[1:]

	cmd := exec.Command(bin, args...)
	// stdoutStderr, err := cmd.CombinedOutput()
	return cmd.CombinedOutput()
}
