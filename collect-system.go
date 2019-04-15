package ecediag

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/go-sysinfo"
	"github.com/shirou/gopsutil/cpu"
)

func runSystemCmds(tar *Tarball) {
	// l := logp.NewLogger("system_cmd")

	fmt.Println("[ ] Collecting system information")

	var wg sync.WaitGroup
	for _, cmd := range SystemCmd {
		wg.Add(1)
		go cmd.processTask(tar, &wg)
	}
	for _, sf := range SystemFiles {
		wg.Add(1)
		go sf.processTask(tar, &wg)
	}
	wg.Wait()

	// fpath := filepath.Join(cfg.DiagName, "server_info", c.Filename)
	fp := func(path string) string { return filepath.Join(cfg.DiagName, "server_info", path) }

	sysInfo := sysinfo.Go()
	writeJSON(fp("GoSysInfo.txt"), sysInfo, tar)
	hostInfo, _ := sysinfo.Host()
	writeJSON(fp("GoHostInfo.txt"), hostInfo, tar)
	procs, _ := sysinfo.Processes()

	procInfo := make([]interface{}, 100)
	for _, proc := range procs {
		info, _ := proc.Info()
		procInfo = append(procInfo, info)
		// if err != nil {
		// 	if os.IsPermission(err) {
		// 		continue
		// 	}
		// 	t.Fatal(err)
		// }
	}
	writeJSON(fp("GoProcInfo.txt"), procInfo, tar)

	cpuInfo, _ := cpu.Info()
	writeJSON(fp("GoCPUinfo.txt"), cpuInfo, tar)
	cpuTimeStat, _ := cpu.Times(false)
	writeJSON(fp("GoCPUtimeStat.txt"), cpuTimeStat, tar)

	// resc, errc := make(chan string), make(chan error)
	// for _, cmd := range SystemCmd {
	// 	go func(cmd interface{}) {
	// 		body, err := processTask(cmd, tar)
	// 		// body, err := runCommand(cmd)
	// 		if err != nil {
	// 			errc <- err
	// 			return
	// 		}
	// 		resc <- string(body)
	// 	}(cmd)
	// }

	// for i := 0; i < len(SystemCmd); i++ {
	// 	select {
	// 	case res := <-resc:
	// 		l.Info(res)
	// 	case err := <-errc:
	// 		l.Error(err)
	// 	}
	// }
	clearStdoutLine()
	fmt.Println("\r[✔] Collected system information")
}

func (c systemFile) processTask(tar *Tarball, wg *sync.WaitGroup) {
	l := logp.NewLogger("system_file")
	files, err := c.checkFile()
	if err != nil {
		l.Error(err)
	} else {
		fp := func(path string) string { return filepath.Join(cfg.DiagName, "server_info", path) }

		var buf []byte
		numFiles := len(files)

		if numFiles == 1 {
			stat, _ := os.Stat(files[0])
			tar.AddFile(files[0], stat, fp(c.Filename))
			l.Infof("Collected %s as %s", files[0], c.Filename)
		} else {
			for _, file := range files {
				fileData, _ := ioutil.ReadFile(file)
				header := fmt.Sprintf("==> %s <==\n", file)
				buf = append(buf, []byte(header)...)
				buf = append(buf, fileData...)
				buf = append(buf, []byte("\n")...)
			}
			tar.AddData(fp(c.Filename), buf)
			l.Infof("Combined contents of %v into %s", files, c.Filename)
		}
	}
	wg.Done()
}

func (c systemCmd) processTask(tar *Tarball, wg *sync.WaitGroup) {
	l := logp.NewLogger("system_cmd")
	out, err := c.executeCmd()
	if err != nil {
		l.Error(err)
	} else {
		fpath := filepath.Join(cfg.DiagName, "server_info", c.Filename)
		tar.AddData(fpath, out)
		l.Infof("Command completed: \"%v\" -> %s", c.RawCmd, c.Filename)
	}
	wg.Done()
}

func (c *systemCmd) executeCmd() ([]byte, error) {
	output, err := c.run()
	if err != nil {
		err = fmt.Errorf("Command failed: %s, `%v`, %s", c.Filename, c.RawCmd, err)
		// return output, err
	}
	return output, err
}

func (c *systemCmd) run() ([]byte, error) {
	expand := strings.Split(c.RawCmd, " ")
	bin := expand[0]
	args := expand[1:]

	cmd := exec.Command(bin, args...)
	// stdoutStderr, err := cmd.CombinedOutput()
	return cmd.CombinedOutput()
}

func (c *systemFile) checkFile() ([]string, error) {
	files, err := filepath.Glob(c.RawFile)
	if err != nil {
		return nil, err
	}
	if len(files) > 0 {
		return files, nil
	}
	return nil, fmt.Errorf("No files found for pattern %s", c.RawFile)
}
