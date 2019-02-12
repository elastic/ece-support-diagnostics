package ece_support_diagnostic

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/elastic/beats/libbeat/logp"
)

type command struct {
	filename string
	cmd      string
}

func SystemCommands() {
	l := logp.NewLogger("system_cmd")
	commands := []command{
		command{"uname.txt", "uname -a"},
		command{"top.txt", "top -b -n1"},
		command{"ps.txt", "ps -eaf"},
		command{"df.txt", "df -h"},
		command{"fs_permissions_storage_path.txt", "ls -al $storage_path"},
		command{"fs_permissions_mnt_data.txt", "ls -al /mnt/data"},

		command{"netstat_all.txt", "netstat -anp"},
		command{"netstat_listening.txt", "netstat -ntulpn"},
		command{"iptables.txt", "iptables -L"},
		command{"routes.txt", "route -n"},
		command{"mounts.txt", "mount"},

		// command{"netstat_all.txt", "sudo netstat -anp"},
		// command{"netstat_listening.txt", "sudo netstat -ntulpn"},
		// command{"iptables.txt", "sudo iptables -L"},
		// command{"routes.txt", "sudo route -n"},
		// command{"mounts.txt", "sudo mount"},

		command{"linux-release.txt", "cat /etc/*-release"},
		command{"fstab.txt", "cat /etc/fstab"},
		// command{"fstab.txt", "sudo cat /etc/fstab"},
		command{"cpu.txt", "cat /proc/cpuinfo"},
		command{"limits.txt", "cat /etc/security/limits.conf"},

		command{"top_threads.txt", "top -b -n1 -H"},
		command{"sysctl.txt", "sysctl -a"},
		command{"dmesg.txt", "dmesg"},
		command{"ss.txt", "ss"},
		command{"iostat.txt", "iostat -c -d -x -t -m 1 5"},
	}
	// if [ -x "$(type -P sar)" ];
	//   then
	//     #sar individual devices - sample 5 times every 1 second
	//     print_msg "SAR [sampling individual I/O devices]" "INFO"
	//     sar -d -p 1 5 > $elastic_folder/sar_devices.txt 2>&1
	//     #CPU usage - individual cores - sample 5 times every 1 second
	//     print_msg "SAR [sampling CPU cores usage]" "INFO"
	//     sar -P ALL 1 5 > $elastic_folder/sar_cpu_cores.txt 2>&1
	//     #load average last 1-5-15 minutes - 1 sample
	//     print_msg "SAR [collect load average]" "INFO"
	//     sar -q 1 1 > $elastic_folder/sar_load_average_sampled.txt 2>&1
	//     #memory - sample 5 times every 1 second
	//     print_msg "SAR [sampling memory usage]" "INFO"
	//     sar -r 1 5 > $elastic_folder/sar_memory_sampled.txt 2>&1
	//     #swap - sample once
	//     print_msg "SAR [collect swap usage]" "INFO"
	//     sar -S 1 1 > $elastic_folder/sar_swap_sampled.txt 2>&1
	//     #network
	//     print_msg "SAR [collect network stats]" "INFO"
	//     sar -n DEV > $elastic_folder/sar_network.txt 2>&1
	//   else
	//     print_msg "'sar' command not found. Please install package 'sysstat' to collect extended system stats" "WARN"

	// print_msg "Grabbing ECE logs" "INFO"
	// cd $storage_path && find . -type f -name *.log -exec cp --parents \{\} $elastic_folder \;
	// print_msg "Checking XFS info" "INFO"
	// [[ -x "$(type -P xfs_info)" ]] && xfs_info $storage_path > $elastic_folder/xfs_info.txt 2>&1

	resc, errc := make(chan string), make(chan error)
	for _, cmd := range commands {
		go func(cmd command) {
			body, err := runCommand(cmd)
			if err != nil {
				errc <- err
				return
			}
			resc <- string(body)
		}(cmd)
	}

	for i := 0; i < len(commands); i++ {
		select {
		case res := <-resc:
			l.Info(res)
		case err := <-errc:
			l.Error(err)
		}
	}
}

func runCommand(item command) (string, error) {
	// print_msg "Gathering system info..." "INFO"
	ex := strings.Split(item.cmd, " ")
	bin := ex[0]
	args := ex[1:]

	cmd := exec.Command(bin, args...)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		err = errors.New(fmt.Sprintf("Command failed: `%v`, %s", item.cmd, err))
		return "", err
	}

	// Write data to file
	filepath := filepath.Join(SystemTmpDir, item.filename)
	err = ioutil.WriteFile(filepath, stdoutStderr, 0644)
	if err != nil {
		return "", err
	}
	res := fmt.Sprintf("Command completed: \"%v\"", item.cmd)
	return res, nil
}
