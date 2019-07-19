package systemInfo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/ece-support-diagnostics/config"
	"github.com/elastic/ece-support-diagnostics/helpers"
	"github.com/elastic/go-sysinfo"
	"github.com/shirou/gopsutil/cpu"
)

type systemInfo struct{}

// Run runs
func Run(sCmd []SystemCmd, sFiles []SystemFile, cfg *config.Config) {
	sysinfo := systemInfo{}
	sysinfo.runSystemCmds(sCmd, sFiles, cfg)
}

func (s systemInfo) runSystemCmds(sCmd []SystemCmd, sFiles []SystemFile, cfg *config.Config) {
	// l := logp.NewLogger("system_cmd")

	fmt.Println("[ ] Collecting system information")

	var wg sync.WaitGroup
	for _, cmd := range sCmd {
		wg.Add(1)
		go s.processCmd(cmd, &wg, cfg)
	}
	for _, sf := range sFiles {
		wg.Add(1)
		go s.processFile(sf, &wg, cfg)
	}
	wg.Wait()

	fp := func(path string) string { return filepath.Join(cfg.DiagnosticFilename(), "server_info", path) }

	sysInfo := sysinfo.Go()
	s.writeJSON(fp("GoSysInfo.txt"), sysInfo, cfg)
	// writeJSON(fp("GoSysInfo.txt"), sysInfo, tar)

	hostInfo, _ := sysinfo.Host()
	s.writeJSON(fp("GoHostInfo.txt"), hostInfo, cfg)
	// writeJSON(fp("GoHostInfo.txt"), hostInfo, tar)

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

	s.writeJSON(fp("GoProcInfo.txt"), procInfo, cfg)
	// writeJSON(fp("GoProcInfo.txt"), procInfo, tar)

	cpuInfo, _ := cpu.Info()
	s.writeJSON(fp("GoCPUinfo.txt"), cpuInfo, cfg)
	// writeJSON(fp("GoCPUinfo.txt"), cpuInfo, tar)

	cpuTimeStat, _ := cpu.Times(false)
	s.writeJSON(fp("GoCPUtimeStat.txt"), cpuTimeStat, cfg)

	// writeJSON(fp("GoCPUtimeStat.txt"), cpuTimeStat, tar)

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
	helpers.ClearStdoutLine()
	fmt.Println("\r[✔] Collected system information")
}

func (s systemInfo) processFile(c SystemFile, wg *sync.WaitGroup, cfg *config.Config) {
	l := logp.NewLogger("system_file")
	files, err := s.checkFile(&c)
	if err != nil {
		l.Error(err)
	} else {
		fp := func(path string) string { return filepath.Join(cfg.DiagnosticFilename(), "server_info", path) }

		var buf []byte
		numFiles := len(files)

		if numFiles == 1 {
			stat, _ := os.Stat(files[0])
			cfg.Store.AddFile(files[0], stat, fp(c.Filename))
			l.Infof("Collected %s as %s", files[0], c.Filename)
		} else {
			for _, file := range files {
				fileData, _ := ioutil.ReadFile(file)
				header := fmt.Sprintf("==> %s <==\n", file)
				buf = append(buf, []byte(header)...)
				buf = append(buf, fileData...)
				buf = append(buf, []byte("\n")...)
			}
			cfg.Store.AddData(fp(c.Filename), buf)
			l.Infof("Combined contents of %v into %s", files, c.Filename)
		}
	}
	wg.Done()
}

func (s systemInfo) processCmd(c SystemCmd, wg *sync.WaitGroup, cfg *config.Config) {
	l := logp.NewLogger("system_cmd")
	out, err := s.executeCmd(&c)
	if err != nil {
		l.Error(err)
	} else {
		fpath := filepath.Join(cfg.DiagnosticFilename(), "server_info", c.Filename)
		cfg.Store.AddData(fpath, out)
		l.Infof("Command completed: \"%v\" -> %s", c.RawCmd, c.Filename)
	}
	wg.Done()
}

func (s systemInfo) executeCmd(c *SystemCmd) ([]byte, error) {
	output, err := s.run(c)
	if err != nil {
		err = fmt.Errorf("Command failed: %s, `%v`, %s", c.Filename, c.RawCmd, err)
		// return output, err
	}
	return output, err
}

func (s systemInfo) run(c *SystemCmd) ([]byte, error) {
	expand := strings.Split(c.RawCmd, " ")
	bin := expand[0]
	args := expand[1:]

	cmd := exec.Command(bin, args...)
	// stdoutStderr, err := cmd.CombinedOutput()
	return cmd.CombinedOutput()
}

func (s systemInfo) checkFile(c *SystemFile) ([]string, error) {
	files, err := filepath.Glob(c.RawFile)
	if err != nil {
		return nil, err
	}
	if len(files) > 0 {
		return files, nil
	}
	return nil, fmt.Errorf("No files found for pattern %s", c.RawFile)
}

func (s systemInfo) writeJSON(path string, apiResp interface{}, cfg *config.Config) error {
	json, err := json.MarshalIndent(apiResp, "", "  ")
	if err != nil {
		panic(err)
	}
	err = cfg.Store.AddData(path, json)
	if err != nil {
		panic(err)
	}
	return err
}
