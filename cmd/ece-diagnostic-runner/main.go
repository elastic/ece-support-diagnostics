package main

import (
	"os"

	sd "github.com/elastic/ece-support-diagnostics"
)

// func init() {
// 	if cpu := runtime.NumCPU(); cpu == 1 {
// 		runtime.GOMAXPROCS(2)
// 	} else {
// 		runtime.GOMAXPROCS(cpu)
// 	}
// }

func init() {

}

func main() {
	if err := sd.Start(); err != nil {
		os.Exit(1)
	}
}
