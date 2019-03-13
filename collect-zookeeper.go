package ecediag

import (
	"errors"
	"io"
	"net"
	"os/exec"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/elastic/beats/libbeat/logp"
)

func zookeeperMNTR(c types.Container, tar *Tarball) {
	log := logp.NewLogger("zookeeper")
	log.Info("Collecting zookeeper mntr")

	var port uint16
	for _, p := range c.Ports {
		if p.PublicPort >= 2100 && p.PublicPort <= 2199 {
			port = p.PublicPort
		}
	}
	if port == 0 {
		log.Error("Could not determine Zookeeper port")
		return
	}

	ip, err := externalIP()
	if err != nil {
		log.Errorf("Not collecting `mntr` info, %s", err)
		return
	}

	cmd := exec.Command("nc", ip, "2193")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "mntr")
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	fpath := filepath.Join(DiagName, "zookeeper_mntr.txt")
	tar.AddData(fpath, out)
	// fmt.Println(test, err)
}

func externalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		if ip == nil || ip.IsLoopback() {
			continue
		}
		ip = ip.To4()
		if ip == nil {
			continue // not an ipv4 address
		}
		return ip.String(), nil
	}
	return "", errors.New("Could not determine an ipv4 address")
}
