package main

import (
	"flag"

	sd "github.com/elastic/ece-support-diagnostics"
	"github.com/elastic/ece-support-diagnostics/config"
)

var (
	// Basepath provides the tmp location to create the tar file in
	basepath string
	// ElasticFolder provides the path to where ECE is installed
	elasticfolder string
	// DisableRest is used for disabling collecting Rest/HTTP requests
	disablerest bool
	// UploadUID provides the unique ID that needs to be specified for using the Elastic upload service
	uploadUID string
)

func init() {
	flag.StringVar(&basepath, "t", "/tmp", "Path to the elastic folder")
	flag.StringVar(&elasticfolder, "f", "/mnt/data/elastic", "Path to the elastic folder")
	flag.BoolVar(&disablerest, "disableRest", false, "Disable Rest calls")
	flag.StringVar(&uploadUID, "u", "", "Elastic Upload ID")
	flag.Parse()
}

func main() {
	ece := config.New()
	ece.Basepath = basepath
	ece.ElasticFolder = elasticfolder
	ece.DisableRest = disablerest
	ece.UploadUID = uploadUID
	sd.Start(ece)
	// if err := sd.Start(); err != nil {
	// 	os.Exit(1)
	// }
}
