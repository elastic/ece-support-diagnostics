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
	"github.com/elastic/ece-support-diagnostics/store"
	"github.com/elastic/go-sysinfo"
	"github.com/shirou/gopsutil/cpu"
)

type fileSystemStore struct {
	store.ContentStore
	cfg *config.Config
}

// Run runs
func Run(t store.ContentStore, sCmd []SystemCmd, sFiles []SystemFile, config *config.Config) {
	store := fileSystemStore{t, config}
	store.runSystemCmds(sCmd, sFiles)
}

func (t fileSystemStore) runSystemCmds(sCmd []SystemCmd, sFiles []SystemFile) {
	// l := logp.NewLogger("system_cmd")

	fmt.Println("[ ] Collecting system information")

	var wg sync.WaitGroup
	for _, cmd := range sCmd {
		wg.Add(1)
		go t.processCmd(cmd, &wg)
	}
	for _, sf := range sFiles {
		wg.Add(1)
		go t.processFile(sf, &wg)
	}
	wg.Wait()

	// fpath := filepath.Join(cfg.DiagName, "server_info", c.Filename)
	fp := func(path string) string { return filepath.Join(t.cfg.DiagName, "server_info", path) }

	sysInfo := sysinfo.Go()
	t.writeJSON(fp("GoSysInfo.txt"), sysInfo)
	// writeJSON(fp("GoSysInfo.txt"), sysInfo, tar)

	hostInfo, _ := sysinfo.Host()
	t.writeJSON(fp("GoHostInfo.txt"), hostInfo)
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

	t.writeJSON(fp("GoProcInfo.txt"), procInfo)
	// writeJSON(fp("GoProcInfo.txt"), procInfo, tar)

	cpuInfo, _ := cpu.Info()
	t.writeJSON(fp("GoCPUinfo.txt"), cpuInfo)
	// writeJSON(fp("GoCPUinfo.txt"), cpuInfo, tar)

	cpuTimeStat, _ := cpu.Times(false)
	t.writeJSON(fp("GoCPUtimeStat.txt"), cpuTimeStat)

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

func (t fileSystemStore) processFile(c SystemFile, wg *sync.WaitGroup) {
	l := logp.NewLogger("system_file")
	files, err := t.checkFile(&c)
	if err != nil {
		l.Error(err)
	} else {
		fp := func(path string) string { return filepath.Join(t.cfg.DiagName, "server_info", path) }

		var buf []byte
		numFiles := len(files)

		if numFiles == 1 {
			stat, _ := os.Stat(files[0])
			t.AddFile(files[0], stat, fp(c.Filename))
			l.Infof("Collected %s as %s", files[0], c.Filename)
		} else {
			for _, file := range files {
				fileData, _ := ioutil.ReadFile(file)
				header := fmt.Sprintf("==> %s <==\n", file)
				buf = append(buf, []byte(header)...)
				buf = append(buf, fileData...)
				buf = append(buf, []byte("\n")...)
			}
			t.AddData(fp(c.Filename), buf)
			l.Infof("Combined contents of %v into %s", files, c.Filename)
		}
	}
	wg.Done()
}

func (t fileSystemStore) processCmd(c SystemCmd, wg *sync.WaitGroup) {
	l := logp.NewLogger("system_cmd")
	out, err := t.executeCmd(&c)
	if err != nil {
		l.Error(err)
	} else {
		fpath := filepath.Join(t.cfg.DiagName, "server_info", c.Filename)
		t.AddData(fpath, out)
		l.Infof("Command completed: \"%v\" -> %s", c.RawCmd, c.Filename)
	}
	wg.Done()
}

func (t fileSystemStore) executeCmd(c *SystemCmd) ([]byte, error) {
	output, err := t.run(c)
	if err != nil {
		err = fmt.Errorf("Command failed: %s, `%v`, %s", c.Filename, c.RawCmd, err)
		// return output, err
	}
	return output, err
}

func (t fileSystemStore) run(c *SystemCmd) ([]byte, error) {
	expand := strings.Split(c.RawCmd, " ")
	bin := expand[0]
	args := expand[1:]

	cmd := exec.Command(bin, args...)
	// stdoutStderr, err := cmd.CombinedOutput()
	return cmd.CombinedOutput()
}

func (t fileSystemStore) checkFile(c *SystemFile) ([]string, error) {
	files, err := filepath.Glob(c.RawFile)
	if err != nil {
		return nil, err
	}
	if len(files) > 0 {
		return files, nil
	}
	return nil, fmt.Errorf("No files found for pattern %s", c.RawFile)
}

func (t fileSystemStore) writeJSON(path string, apiResp interface{}) error {
	json, err := json.MarshalIndent(apiResp, "", "  ")
	if err != nil {
		panic(err)
	}
	err = t.AddData(path, json)
	if err != nil {
		panic(err)
	}
	return err
}
