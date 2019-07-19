package systemInfo

func NewSystemFileTasks() []SystemFile {
	// SystemFiles will collect files from the target host
	//  The RawFile runs a glob to find files. If more than one file is found
	//  the data from each file is concatenated into the target Filename
	return []SystemFile{
		SystemFile{
			Filename: "linux-release.txt",
			RawFile:  "/etc/*-release",
		},
		SystemFile{
			Filename: "fstab.txt",
			RawFile:  "/etc/fstab",
		},
		SystemFile{
			Filename: "limits.txt",
			RawFile:  "/etc/security/limits.conf",
		},
	}

}

func NewSystemCmdTasks() []SystemCmd {
	// SystemCmds holds the set of system commands that need to be collected
	return []SystemCmd{
		SystemCmd{
			Filename: "uname.txt",
			RawCmd:   "uname -a",
		},
		SystemCmd{
			Filename: "top.txt",
			RawCmd:   "top -b -n1",
		},
		SystemCmd{
			Filename: "ps.txt",
			RawCmd:   "ps -eaf",
		},
		SystemCmd{
			Filename: "df.txt",
			RawCmd:   "df -h",
		},
		// cfg.ElasticFolder is not initalized at this point, this does not work.
		// systemCmd{
		// 	Filename: "fs_permissions_storage_path.txt",
		// 	RawCmd:   fmt.Sprintf("ls -la %s", cfg.ElasticFolder),
		// },
		SystemCmd{
			Filename: "fs_permissions_mnt_data.txt",
			RawCmd:   "ls -la /mnt/data",
		},

		SystemCmd{
			Filename: "netstat_all.txt",
			RawCmd:   "netstat -anp",
		},
		SystemCmd{
			Filename: "netstat_listening.txt",
			RawCmd:   "netstat -ntulpn",
		},
		SystemCmd{
			Filename: "iptables.txt",
			RawCmd:   "iptables -L",
		},
		SystemCmd{
			Filename: "routes.txt",
			RawCmd:   "route -n",
		},
		SystemCmd{
			Filename: "mounts.txt",
			RawCmd:   "mount",
		},
		SystemCmd{
			Filename: "systemd.unit",
			RawCmd:   "systemctl cat docker.service",
		},

		// command{"netstat_all.txt", "sudo netstat -anp"},
		// command{"netstat_listening.txt", "sudo netstat -ntulpn"},
		// command{"iptables.txt", "sudo iptables -L"},
		// command{"routes.txt", "sudo route -n"},
		// command{"mounts.txt", "sudo mount"},

		SystemCmd{
			Filename: "cpu.txt",
			RawCmd:   "cat /proc/cpuinfo",
		},
		// Not sure if this is needed.
		SystemCmd{
			Filename: "fips.txt",
			RawCmd:   "cat /proc/sys/crypto/fips_enabled",
		},
		SystemCmd{
			Filename: "top_threads.txt",
			RawCmd:   "top -b -n1 -H",
		},
		SystemCmd{
			Filename: "sysctl.txt",
			RawCmd:   "sysctl -a",
		},
		SystemCmd{
			Filename: "dmesg.txt",
			RawCmd:   "dmesg",
		},
		SystemCmd{
			Filename: "ss.txt",
			RawCmd:   "ss",
		},
		SystemCmd{
			Filename: "iostat.txt",
			RawCmd:   "iostat -c -d -x -t -m 1 5",
		},
	}
	// TODO: Add generic call to get Virtualization type. Maybe hostnamectl (requires systemd)
}
