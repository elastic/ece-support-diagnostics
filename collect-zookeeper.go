package ecediag

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os/exec"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/elastic/beats/libbeat/logp"
)

// zookeeperMNTR sends `echo mntr|nc ip port` for zookeeper
//  could not use localhost or 0.0.0.0, the response gets dropped
//  discovers the zookeep docker port between 2100-2199
//  then sends the command to the first ipv4 address found on the host
func zookeeperMNTR(container types.Container, tar *Tarball) {
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

	ip, err := externalIP()
	if err != nil {
		log.Errorf("Not collecting `mntr` info, %s", err)
		return
	}

	portString := fmt.Sprintf("%d", port)
	cmd := exec.Command("nc", ip, portString)
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
		log.Fatalf("It didn't work:\n%s\n%s", err, out)
	}

	fpath := filepath.Join(cfg.DiagName, "zookeeper_mntr.txt")
	tar.AddData(fpath, out)
	// fmt.Println(test, err)
}

// find an ipv4 address to use
// TODO: look into adding additional error handling, and not sure if ipv6 could be used
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
