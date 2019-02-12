package ece_support_diagnostic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/elastic/beats/libbeat/logp"
)

func DockerCommands() {
	dlog := logp.NewLogger("docker")
	dlog.Info("GOING TO GRAB MY DOCKER")

	cli, err := client.NewClientWithOpts()
	cli.NegotiateAPIVersion(context.Background())

	ctx := context.Background()
	Containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	writeJson("Containers.json", cmd(Containers, err))
	writeJson("Repository.json", cmd(cli.ImageList(ctx, types.ImageListOptions{})))
	writeJson("Info.json", cmd(cli.Info(ctx)))
	writeJson("DiskUsage.json", cmd(cli.DiskUsage(ctx)))
	writeJson("ServerVersion.json", cmd(cli.ServerVersion(ctx)))

	for _, container := range Containers {

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		logOptions := types.ContainerLogsOptions{
			Since:      "72h",
			ShowStdout: true,
			ShowStderr: true,
		}

		dlog.Info("Writing logs for container: ", container.ID[:10])
		reader, err := cli.ContainerLogs(ctx, container.ID, logOptions)
		if err != nil {
			log.Fatal(err)
		}

		fileName := safeFilename(container.ID[:10], container.Names[0], container.Image)
		filePath := fp("logs", fileName)
		outFile, err := os.Create(filePath + ".log")
		// handle err
		defer outFile.Close()
		_, err = io.Copy(outFile, reader)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		cTop, err := cli.ContainerTop(ctx, container.ID, []string{})
		cTjson, err := json.MarshalIndent(cTop, "", "  ")
		err = ioutil.WriteFile(filePath+".top", cTjson, 0644)

	}
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

func writeJson(path string, apiResp interface{}) string {
	json, err := json.MarshalIndent(apiResp, "", "  ")
	fmt.Println(fp(path))
	err = ioutil.WriteFile(fp(path), json, 0644)
	if err != nil {
		panic(err)
	}
	return "hello"
}

func fp(filename ...string) string {
	newPaths := filepath.Join(filename...)
	return filepath.Join(DockerTmpDir, newPaths)
}

// Hack to allow calling writeJson directly
func cmd(api interface{}, err error) interface{} {
	return api
}
