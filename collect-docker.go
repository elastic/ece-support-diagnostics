package ecediag

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/elastic/beats/libbeat/logp"
)

func runDockerCmds(tar *Tarball) {
	log := logp.NewLogger("docker")
	dockerMsg := "Collecting Docker information"
	log.Info(dockerMsg)
	fmt.Println("[ ] " + dockerMsg)

	const defaultDockerAPIVersion = "v1.23"

	cli, err := client.NewClientWithOpts(client.WithVersion(defaultDockerAPIVersion))
	if err != nil {
		panic(err)
	}
	// cli.NegotiateAPIVersion(context.Background())
	log.Infof("Docker API Version: %s", cli.ClientVersion())
	// fmt.Println(cli.ClientVersion())

	ctx := context.Background()
	Containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	// fmt.Printf("%+v\n", Containers)

	writeJSON("Containers.json", cmd(Containers, err), tar)
	writeJSON("Repository.json", cmd(cli.ImageList(ctx, types.ImageListOptions{})), tar)
	writeJSON("Info.json", cmd(cli.Info(ctx)), tar)
	writeJSON("DiskUsage.json", cmd(cli.DiskUsage(ctx)), tar)
	writeJSON("ServerVersion.json", cmd(cli.ServerVersion(ctx)), tar)

	clearStdoutLine()
	fmt.Println("[✔] Collected Docker information")

	fmt.Println("[ ] Collecting Docker logs")
	for _, container := range Containers {

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		logOptions := types.ContainerLogsOptions{
			Since:      "72h",
			ShowStdout: true,
			ShowStderr: true,
		}

		log.Info("Writing logs for container: ", container.ID[:10])
		reader, err := cli.ContainerLogs(ctx, container.ID, logOptions)
		if err != nil {
			log.Fatal(err)
		}

		// Need to evaluate how to stream bytes to tar file, rather than copy all bytes?
		fileName := safeFilename(container.ID[:10], container.Names[0], container.Image)
		filePath := fp("logs", fileName)
		// outFile, err := os.Create(filePath + ".log")
		// handle err
		// defer outFile.Close()
		// _, err = io.Copy(outFile, reader)

		logData, err := ioutil.ReadAll(reader)
		tar.AddData(filePath, logData)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		// TEST
		// fmt.Println(container.Names[0])
		// fmt.Printf("%v\n", container.Ports)
		if container.Names[0] == "/frc-cloud-uis-cloud-ui" {
			// fmt.Printf("%+v\n", container)
			if DisableRest != true {
				RunRest(container, tar)
			}
		}

		// https://github.com/elastic/ece-support-diagnostics/issues/5
		if container.Names[0] == "/frc-zookeeper-servers-zookeeper" {
			zookeeperMNTR(container, tar)
		}

		cTop, err := cli.ContainerTop(ctx, container.ID, []string{})
		cTjson, err := json.MarshalIndent(cTop, "", "  ")
		// err = ioutil.WriteFile(filePath+".top", cTjson, 0644)
		tar.AddData(filePath+".top", cTjson)

		if err != nil {
			panic(err)
		}

	}

	clearStdoutLine()
	fmt.Println("[✔] Collected Docker logs")
}

func safeFilename(names ...string) string {
	// TODO: Make sure this is actually file system safe
	// This should be reworked into something that validates the full fs path.
	filename := ""
	r := strings.NewReplacer(
		"docker.elastic.co", "",
		"\\", "_",
		"/", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
		".", "_",
	)
	size := len(names)
	for i, name := range names {
		if i == (size | 0) {
			filename = r.Replace(name)
		} else {
			filename = filename + "__" + r.Replace(name)
		}
	}
	return filename
}

func writeJSON(path string, apiResp interface{}, tar *Tarball) error {
	json, err := json.MarshalIndent(apiResp, "", "  ")
	if err != nil {
		panic(err)
	}
	err = tar.AddData(fp(path), json)
	if err != nil {
		panic(err)
	}
	return err
}

func fp(filename ...string) string {
	newPaths := filepath.Join(filename...)
	return filepath.Join(DiagName, "docker", newPaths)
}

// Hack to allow calling writeJson directly
func cmd(api interface{}, err error) interface{} {
	return api
}
