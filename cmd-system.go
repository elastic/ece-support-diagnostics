package ecediag

import "github.com/elastic/ece-support-diagnostics/collectors/systemInfo"

// SystemFiles will collect files from the target host
//  The RawFile runs a glob to find files. If more than one file is found
//  the data from each file is concatenated into the target Filename
var SystemFiles = []systemInfo.SystemFile{
	systemInfo.SystemFile{
		Filename: "linux-release.txt",
		RawFile:  "/etc/*-release",
	},
	systemInfo.SystemFile{
		Filename: "fstab.txt",
		RawFile:  "/etc/fstab",
	},
	systemInfo.SystemFile{
		Filename: "limits.txt",
		RawFile:  "/etc/security/limits.conf",
	},
}

// SystemCmds holds the set of system commands that need to be collected
var SystemCmds = []systemInfo.SystemCmd{
	systemInfo.SystemCmd{
		Filename: "uname.txt",
		RawCmd:   "uname -a",
	},
	systemInfo.SystemCmd{
		Filename: "top.txt",
		RawCmd:   "top -b -n1",
	},
	systemInfo.SystemCmd{
		Filename: "ps.txt",
		RawCmd:   "ps -eaf",
	},
	systemInfo.SystemCmd{
		Filename: "df.txt",
		RawCmd:   "df -h",
	},
	// cfg.ElasticFolder is not initalized at this point, this does not work.
	// systemCmd{
	// 	Filename: "fs_permissions_storage_path.txt",
	// 	RawCmd:   fmt.Sprintf("ls -la %s", cfg.ElasticFolder),
	// },
	systemInfo.SystemCmd{
		Filename: "fs_permissions_mnt_data.txt",
		RawCmd:   "ls -la /mnt/data",
	},

	systemInfo.SystemCmd{
		Filename: "netstat_all.txt",
		RawCmd:   "netstat -anp",
	},
	systemInfo.SystemCmd{
		Filename: "netstat_listening.txt",
		RawCmd:   "netstat -ntulpn",
	},
	systemInfo.SystemCmd{
		Filename: "iptables.txt",
		RawCmd:   "iptables -L",
	},
	systemInfo.SystemCmd{
		Filename: "routes.txt",
		RawCmd:   "route -n",
	},
	systemInfo.SystemCmd{
		Filename: "mounts.txt",
		RawCmd:   "mount",
	},
	systemInfo.SystemCmd{
		Filename: "systemd.unit",
		RawCmd:   "systemctl cat docker.service",
	},

	// command{"netstat_all.txt", "sudo netstat -anp"},
	// command{"netstat_listening.txt", "sudo netstat -ntulpn"},
	// command{"iptables.txt", "sudo iptables -L"},
	// command{"routes.txt", "sudo route -n"},
	// command{"mounts.txt", "sudo mount"},

	systemInfo.SystemCmd{
		Filename: "cpu.txt",
		RawCmd:   "cat /proc/cpuinfo",
	},
	// Not sure if this is needed.
	systemInfo.SystemCmd{
		Filename: "fips.txt",
		RawCmd:   "cat /proc/sys/crypto/fips_enabled",
	},
	systemInfo.SystemCmd{
		Filename: "top_threads.txt",
		RawCmd:   "top -b -n1 -H",
	},
	systemInfo.SystemCmd{
		Filename: "sysctl.txt",
		RawCmd:   "sysctl -a",
	},
	systemInfo.SystemCmd{
		Filename: "dmesg.txt",
		RawCmd:   "dmesg",
	},
	systemInfo.SystemCmd{
		Filename: "ss.txt",
		RawCmd:   "ss",
	},
	systemInfo.SystemCmd{
		Filename: "iostat.txt",
		RawCmd:   "iostat -c -d -x -t -m 1 5",
	},
}

// TODO: Add generic call to get Virtualization type. Maybe hostnamectl (requires systemd)
