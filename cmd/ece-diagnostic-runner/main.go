package main

import (
	"os"

	"github.com/elastic/beats/libbeat/logp"
	sd "github.com/elastic/ece-support-diagnostics"
)

// func init() {
// 	if cpu := runtime.NumCPU(); cpu == 1 {
// 		runtime.GOMAXPROCS(2)
// 	} else {
// 		runtime.GOMAXPROCS(cpu)
// 	}
// }

// func main() {
// 	pt := pt.PlatinumSearcher{Out: os.Stdout, Err: os.Stderr}
// 	exitCode := pt.Run(os.Args[1:])
// 	os.Exit(exitCode)
// }

// var ElasticFolder string

func init() {
	// flag.StringVar(&ElasticFolder, "f", "/mnt/data/elastic", "Path to the elastic folder")
	config := logp.DefaultConfig()
	config.Level = 8
	config.Beat = "HELLO WORLD!"
	config.ToStderr = true
	logp.Configure(config)
}

func main() {
	if err := sd.Start(); err != nil {
		os.Exit(1)
	}
}
