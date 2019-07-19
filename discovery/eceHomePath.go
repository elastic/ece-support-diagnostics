package discovery

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

// CheckStoragePath checks the filesystem for a known folder structure for an ECE install
func CheckStoragePath(eceInstallPath string) (string, error) {

	// ContainerSets stores the roles/containers detected from the runners->runner filesystem folders
	//  this is not currently used for anything
	var ContainerSets = []string{}

	// TODO: if no explicit command line flag overrides, should try looking at docker API
	//  to determine if /mnt/data/elastic is the correct path.

	sp := filepath.Join(eceInstallPath, "*/services/runners/containers/docker/")
	fmt.Printf("Checking the ECE install path, %s\n", sp)
	installPaths, err := filepath.Glob(sp)
	if err != nil {
		log.Fatal(err)
	}

	if len(installPaths) == 1 {
		install := installPaths[0]
		splitPath := strings.Split(strings.TrimRight(install, "/"), "/")
		runnerName := splitPath[len(splitPath)-5]

		roledirs, err := ioutil.ReadDir(install)
		if err != nil {
			log.Fatal(err)
		}

		for _, it := range roledirs {
			roleDir := filepath.Join(install, it.Name())
			fileContainers, err := ioutil.ReadDir(roleDir)
			if err != nil {
				log.Fatal(err)
			}
			for _, ContainerSet := range fileContainers {
				ContainerSets = append(ContainerSets, ContainerSet.Name())
			}
		}
		fmt.Printf("Discovered: %v\n", ContainerSets)
		return runnerName, nil
	}
	return "", fmt.Errorf("could not find a valid ECE install location")
}

// /mnt/data/elastic/172.16.0.10/services/
// runners/
// └── containers
//     └── docker
//         ├── admin-consoles
//         │   └── admin-console
//         ├── allocator-metricbeats
//         │   └── allocator-metricbeat
//         ├── allocators
//         │   └── allocator
//         ├── beats-runners
//         │   └── beats-runner
//         ├── blueprints
//         │   └── blueprint
//         ├── client-forwarders
//         │   └── client-forwarder
//         ├── cloud-uis
//         │   └── cloud-ui
//         ├── constructors
//         │   └── constructor
//         ├── curators
//         │   └── curator
//         ├── directors
//         │   └── director
//         ├── proxies
//         │   └── proxy
//         ├── runners
//         │   └── runner
//         ├── services-forwarders
//         │   └── services-forwarder
//         └── zookeeper-servers
//             └── zookeeper
