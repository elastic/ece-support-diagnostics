package ecediag

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/elastic/beats/libbeat/logp"
)

var re = regexp.MustCompile(`\/f[r|a]c-(\w+(?:-\w+)?)-(\w+)`)

func runDockerCmds(tar *Tarball) {
	// log := logp.NewLogger("docker")
	dockerMsg := "Collecting Docker information"
	logp.Info(dockerMsg)
	fmt.Println("[ ] " + dockerMsg)

	const defaultDockerAPIVersion = "v1.23"

	cli, err := client.NewClientWithOpts(client.WithVersion(defaultDockerAPIVersion))
	if err != nil {
		panic(err)
	}
	// cli.NegotiateAPIVersion(context.Background())
	logp.Info("Docker API Version: %s", cli.ClientVersion())
	// fmt.Println(cli.ClientVersion())

	ctx := context.Background()
	Containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	// fmt.Printf("%+v\n", Containers)

	fp := func(path string) string { return filepath.Join(cfg.DiagName, path) }
	writeJSON(fp("DockerContainers.json"), cmd(Containers, err), tar)

	writeJSON(fp("DockerRepository.json"), cmd(cli.ImageList(ctx, types.ImageListOptions{})), tar)

	writeJSON(fp("DockerInfo.json"), cmd(cli.Info(ctx)), tar)

	writeJSON(fp("DockerDiskUsage.json"), cmd(cli.DiskUsage(ctx)), tar)

	writeJSON(fp("DockerServerVersion.json"), cmd(cli.ServerVersion(ctx)), tar)

	clearStdoutLine()
	fmt.Println("[✔] Collected Docker information")

	for _, container := range Containers {
		if container.Names[0] == "/frc-cloud-uis-cloud-ui" {
			if cfg.DisableRest != true {
				RunRest(container, tar)
			}
		}

		// https://github.com/elastic/ece-support-diagnostics/issues/5
		if container.Names[0] == "/frc-zookeeper-servers-zookeeper" {
			zookeeperMNTR(container, tar)
			fmt.Println("[✔] Collected Zookeeper data")
		}

		fmt.Println("[ ] Collecting Docker logs")

		dockerLogs(cli, container, tar)

		clearStdoutLine()
		fmt.Println("[✔] Collected Docker logs")
	}

}

func dockerLogs(cli *client.Client, container types.Container, tar *Tarball) {

	// getValue := func(key string) string { if val, ok := container.Labels[key]; ok {return val}; return "" }
	dockerName := container.Names[0]

	if strings.HasPrefix(dockerName, "/frc") || strings.HasPrefix(dockerName, "/fac") {
		filePath := createFilePath(container)
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		logOptions := types.ContainerLogsOptions{
			Since:      "72h",
			ShowStdout: true,
			ShowStderr: true,
		}

		logp.Info("Writing logs for container: %s", container.ID[:10])
		reader, err := cli.ContainerLogs(ctx, container.ID, logOptions)
		if err != nil {
			logp.Error(err)
		}

		// Need to evaluate how to stream bytes to tar file, rather than copy all bytes?
		// fileName := safeFilename(container.ID[:10], container.Names[0], container.Image)
		// filePath := fp("logs", fileName)
		// outFile, err := os.Create(filePath + ".log")
		// handle err
		// defer outFile.Close()
		// _, err = io.Copy(outFile, reader)

		logData, err := ioutil.ReadAll(reader)
		tar.AddData(filePath+".log", logData)
		if err != nil && err != io.EOF {
			logp.Error(err)
		}

		cTop, err := cli.ContainerTop(ctx, container.ID, []string{})
		cTjson, err := json.MarshalIndent(cTop, "", "  ")
		// err = ioutil.WriteFile(filePath+".top", cTjson, 0644)
		tar.AddData(filePath+".top", cTjson)

		if err != nil {
			panic(err)
		}
	}

}

// TODO: FIX THIS, regex should be passed in, not global
// func createFilePath(container types.Container, re *regexp.Regexp) string {
func createFilePath(container types.Container) string {
	dockerName := container.Names[0]
	labels := container.Labels

	eceLogPath := func(kind, filename string) string {
		return filepath.Join(cfg.DiagName, "ece", kind, filename)
	}
	// if a runner launches the container then it has `runner.container_name`
	if containerName, ok := labels["co.elastic.cloud.runner.container_name"]; ok {
		version := strings.Split(labels["org.label-schema.version"], "-")[0] // "2.1.0-SNAPSHOT"
		fileName := fmt.Sprintf("%s_%s_%.12s", containerName, version, container.ID)

		return eceLogPath(containerName, fileName)
		// cloud-ui_2.1.0_0b3a6993d552.log
	}

	// if an allocator launches the container then it as `allocator.kind`
	if kind, ok := labels["co.elastic.cloud.allocator.kind"]; ok {
		// "elasticsearch" | "kibana"
		clusterID := labels["co.elastic.cloud.allocator.cluster_id"]     // "c5900a8affb44d108ebe31513480a9b8"
		version := labels["co.elastic.cloud.allocator.type_version"]     // "6.6.0"
		instanceName := labels["co.elastic.cloud.allocator.instance_id"] // "instance-0000000000"
		fileName := fmt.Sprintf("%.6s_%s-%s_%s_%.12s", clusterID, kind, version, instanceName, container.ID+".log")

		return filepath.Join(cfg.DiagName, kind, clusterID, instanceName, fileName)
		// 5a4f7f_elasticsearch-5.6.14_instance-0000000000_506b8c016045.log
	}

	// these should be special containers that auto start themselves on reboot
	//  thus they do not have any of the docker Labels above
	//  this also serves as a catch all
	var name string
	match := re.FindStringSubmatch(dockerName)
	if len(match) == 3 {
		name = match[2]
	} else {
		name = strings.TrimPrefix(dockerName, "/frc-")
	}
	version := strings.Split(container.Labels["org.label-schema.version"], "-")[0]
	fileName := fmt.Sprintf("%s_%s_%.12s", name, version, container.ID)
	return eceLogPath(name, fileName)

}

// func safeFilename(names ...string) string {
// 	// TODO: Make sure this is actually file system safe
// 	// This should be reworked into something that validates the full fs path.
// 	filename := ""
// 	r := strings.NewReplacer(
// 		"docker.elastic.co", "",
// 		"\\", "_",
// 		"/", "_",
// 		":", "_",
// 		"*", "_",
// 		"?", "_",
// 		"\"", "_",
// 		"<", "_",
// 		">", "_",
// 		"|", "_",
// 		".", "_",
// 	)
// 	size := len(names)
// 	for i, name := range names {
// 		if i == size || i == 0 {
// 			filename = r.Replace(name)
// 		} else {
// 			filename = filename + "__" + r.Replace(name)
// 		}
// 	}
// 	return filename
// }

func writeJSON(path string, apiResp interface{}, tar *Tarball) error {
	json, err := json.MarshalIndent(apiResp, "", "  ")
	if err != nil {
		panic(err)
	}
	err = tar.AddData(path, json)
	if err != nil {
		panic(err)
	}
	return err
}

// func fp(filename ...string) string {
// 	newPaths := filepath.Join(filename...)
// 	return filepath.Join(cfg.DiagName, "docker", newPaths)
// }

// Hack to allow calling writeJson directly
func cmd(api interface{}, err error) interface{} {
	return api
}
