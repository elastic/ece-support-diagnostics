package zookeeper

import (
	"fmt"
	"io/ioutil"
	"net"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/ece-support-diagnostics/config"
)

type zooCollector struct{}

// Run runs
func Run(c types.Container, cfg *config.Config) {
	zk := zooCollector{}
	zk.zookeeperMNTR(c, cfg)
}

// zookeeperMNTR sends `echo mntr|nc ip port` for zookeeper
//  could not use localhost or 0.0.0.0, the response gets dropped
//  discovers the zookeep docker port between 2100-2199
//  then sends the command to ipv4 address of the docker gateway
func (zk zooCollector) zookeeperMNTR(container types.Container, cfg *config.Config) {
	log := logp.NewLogger("zookeeper")
	log.Info("Collecting zookeeper mntr")

	var port uint16
	for _, p := range container.Ports {
		if p.PublicPort >= 2100 && p.PublicPort <= 2199 {
			port = p.PublicPort
		}
	}
	if port == 0 {
		log.Error("Could not determine Zookeeper port")
		return
	}

	ip := container.NetworkSettings.Networks["bridge"].Gateway

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), 8*time.Second)
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	defer conn.Close()
	fmt.Fprintf(conn, "mntr\n")
	resp, _ := ioutil.ReadAll(conn)

	fpath := filepath.Join(cfg.DiagnosticFilename(), "ece", "zookeeper_mntr.txt")
	cfg.Store.AddData(fpath, resp)

}
