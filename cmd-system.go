package ecediag

import "fmt"

type systemCmd struct {
	Filename string
	RawCmd   string
}
type systemFile struct {
	Filename string
	RawFile  string
}

// SystemFiles will collect files from the target host
//  The RawFile runs a glob to find files. If more than one file is found
//  the data from each file is concatenated into the target Filename
var SystemFiles = []systemFile{
	systemFile{
		Filename: "linux-release.txt",
		RawFile:  "/etc/*-release",
	},
	systemFile{
		Filename: "fstab.txt",
		RawFile:  "/etc/fstab",
	},
	systemFile{
		Filename: "limits.txt",
		RawFile:  "/etc/security/limits.conf",
	},
}

// SystemCmd holds the set of system commands that need to be collected
var SystemCmd = []systemCmd{
	systemCmd{
		Filename: "uname.txt",
		RawCmd:   "uname -a",
	},
	systemCmd{
		Filename: "top.txt",
		RawCmd:   "top -b -n1",
	},
	systemCmd{
		Filename: "ps.txt",
		RawCmd:   "ps -eaf",
	},
	systemCmd{
		Filename: "df.txt",
		RawCmd:   "df -h",
	},
	// cfg.ElasticFolder is not initalized at this point, this does not work.
	systemCmd{
		Filename: "fs_permissions_storage_path.txt",
		RawCmd:   fmt.Sprintf("ls -la %s", cfg.ElasticFolder),
	},
	systemCmd{
		Filename: "fs_permissions_mnt_data.txt",
		RawCmd:   "ls -la /mnt/data",
	},

	systemCmd{
		Filename: "netstat_all.txt",
		RawCmd:   "netstat -anp",
	},
	systemCmd{
		Filename: "netstat_listening.txt",
		RawCmd:   "netstat -ntulpn",
	},
	systemCmd{
		Filename: "iptables.txt",
		RawCmd:   "iptables -L",
	},
	systemCmd{
		Filename: "routes.txt",
		RawCmd:   "route -n",
	},
	systemCmd{
		Filename: "mounts.txt",
		RawCmd:   "mount",
	},
	systemCmd{
		Filename: "systemd.unit",
		RawCmd:   "systemctl cat docker.service",
	},

	// command{"netstat_all.txt", "sudo netstat -anp"},
	// command{"netstat_listening.txt", "sudo netstat -ntulpn"},
	// command{"iptables.txt", "sudo iptables -L"},
	// command{"routes.txt", "sudo route -n"},
	// command{"mounts.txt", "sudo mount"},

	systemCmd{
		Filename: "cpu.txt",
		RawCmd:   "cat /proc/cpuinfo",
	},
	// Not sure if this is needed.
	systemCmd{
		Filename: "fips.txt",
		RawCmd:   "cat /proc/sys/crypto/fips_enabled",
	},
	systemCmd{
		Filename: "top_threads.txt",
		RawCmd:   "top -b -n1 -H",
	},
	systemCmd{
		Filename: "sysctl.txt",
		RawCmd:   "sysctl -a",
	},
	systemCmd{
		Filename: "dmesg.txt",
		RawCmd:   "dmesg",
	},
	systemCmd{
		Filename: "ss.txt",
		RawCmd:   "ss",
	},
	systemCmd{
		Filename: "iostat.txt",
		RawCmd:   "iostat -c -d -x -t -m 1 5",
	},
}

// TODO: Add generic call to get Virtualization type. Maybe hostnamectl (requires systemd)

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
