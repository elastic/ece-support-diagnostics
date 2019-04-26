package checker

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

// ContainerSets is used for storing roles/containers detected from the runners->runner filesystem folders
var ContainerSets = []string{}

// CheckStoragePath checks the filesystem for a known folder structure for an ECE install
func CheckStoragePath(eceInstallPath string) (string, error) {
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
	return "", fmt.Errorf("Could not find a valid ECE install location")
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

// // Read info from JSON. Decided to go with reading file structure on disk instead
// type RunnerData struct {
// 	Containers struct {
// 		Docker map[string]interface{} `json:"docker"`
// 	} `json:"containers"`
// }
//
// type Runners struct {
// 	Containers []Container `json:"Containers"`
// }
//
// type Container struct {
// 	ID struct {
// 		Kind             string `json:"kind"`
// 		ContainerSetName string `json:"container_set_name"`
// 		ContainerName    string `json:"container_name"`
// 		DockerName       string `json:"docker_name"`
// 	} `json:"id"`
// 	Options struct {
// 		Enabled      bool `json:"enabled"`
// 		RolesManaged bool `json:"roles_managed"`
// 		ActiveRoles  map[string]interface {
// 		} `json:"active_roles"`
// 	} `json:"options"`
// }
//
// func (runner *Runners) addItem(it Container) []Container {
// 	runner.Containers = append(runner.Containers, it)
// 	return runner.Containers
// }
//
// func checkStoragePath() {
//
// 	config := &mapstructure.DecoderConfig{
// 		TagName: "json",
// 	}
//
// 	Rn := Runners{}
// 	// path := "/mnt/data/elastic/"
// 	runner := filepath.Join(Basepath, "test_data", "elastic/*/services/runner/managed/container-ids.json")
// 	runnerGlob, err := filepath.Glob(runner)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(runnerGlob)
//
// 	jsonFile, err := os.Open(runnerGlob[0])
// 	// if we os.Open returns an error then handle it
// 	if err != nil {
// 		fmt.Println(err)
// 	}
//
// 	// fmt.Println("Successfully Opened users.json")
// 	defer jsonFile.Close()
//
// 	jsonBytes, _ := ioutil.ReadAll(jsonFile)
// 	var runnerData RunnerData
// 	json.Unmarshal(jsonBytes, &runnerData)
//
// 	for _, containerSetNameMap := range runnerData.Containers.Docker {
// 		for _, containerNameMap := range containerSetNameMap.(map[string]interface{}) {
// 			// fmt.Println(containerName)
// 			cn := Container{}
// 			config.Result = &cn
// 			decoder, _ := mapstructure.NewDecoder(config)
// 			decoder.Decode(containerNameMap)
// 			cn.ID.DockerName = "frc-" + cn.ID.ContainerSetName + "-" + cn.ID.ContainerName
//
// 			Rn.addItem(cn)
// 			fmt.Printf("%+v\n", cn)
// 		}
// 	}
//
// 	// fmt.Printf("%+v\n", Rn)
//
// 	b, err := json.MarshalIndent(Rn, "", " ")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
//
// 	fmt.Println(string(b))
// }
