package eceAPI

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/Masterminds/semver"
	"github.com/elastic/ece-support-diagnostics/collectors/elasticsearch"
	"github.com/elastic/ece-support-diagnostics/config"
	"github.com/elastic/ece-support-diagnostics/helpers"
	"github.com/tidwall/gjson"
)

func RunECEapis(cfg *config.Config) {
	fmt.Println("[ ] Collecting API information ECE and Elasticsearch")

	// store := new(fileStore)

	// Used for templating and version checking.
	ece := helpers.TaskContext{
		// Version: run ECEversionCheck to get the initial version
		Endpoint:  cfg.APIendpoint,
		User:      cfg.User,
		Pass:      cfg.Pass,
		Client:    cfg.HTTPclient,
		Store:     cfg.Store,
		StorePath: cfg.DiagnosticFilename(),
		// Meta:   store a generic interface used for templating variables
		//           in the filename / uri's
	}
	ece.Version = ECEversionCheck(ece)
	// fmt.Println(ece.Version)

	// loop over each Task and dispatch
	func(tasks helpers.Tasks, tskCtl helpers.TaskContext) {

		var wg sync.WaitGroup

		for _, task := range tasks {
			// Add task to wait group
			wg.Add(1)

			ece.Task = task

			go ece.TaskExecuteWithWaitGroup(&wg)
		}

		wg.Wait()

	}(*NewECEtasks(), ece)

	helpers.ClearStdoutLine()
	fmt.Println("[✔] Collected API information for ECE and Elasticsearch")
}

type ECEesClusters struct{}

func (e ECEesClusters) Exec(t helpers.TaskContext, payload []byte) {
	// used to wait for each cluster to complete
	var clusterWait sync.WaitGroup

	esClusters := gjson.GetBytes(payload, "elasticsearch_clusters")

	esClusters.ForEach(func(key, value gjson.Result) bool {
		// clusterName := gjson.Get(value.Raw, "cluster_name")
		clusterID := gjson.Get(value.Raw, "cluster_id")
		versions := gjson.Get(value.Raw, `topology.instances.#.service_version`)

		// println("latest version: ", latestVersion(versions))

		t.Meta = map[string]string{"cluster_id": clusterID.String()}

		// loop over ECE cluster deployment APIs
		go func(tasks helpers.Tasks, tskCtl helpers.TaskContext, clusterWait *sync.WaitGroup) {
			clusterWait.Add(1)
			var wg sync.WaitGroup

			for _, task := range tasks {
				// Add task to wait group
				wg.Add(1)
				tskCtl.Task = task
				go tskCtl.TaskExecuteWithWaitGroup(&wg)
			}

			wg.Wait()
			clusterWait.Done()

		}(*NewECEdeploymentTasks(), t, &clusterWait)

		// println(clusterName.String())
		// println(clusterID.String())
		// println()

		// loop over each Task and dispatch
		go func(tasks helpers.Tasks, tskCtl helpers.TaskContext, clusterWait *sync.WaitGroup) {
			clusterWait.Add(1)
			var wg sync.WaitGroup

			tskCtl.Endpoint = t.Endpoint + "/api/v1/clusters/elasticsearch/" + clusterID.String() + "/proxy/"
			// t.Endpoint = path.Join(t.Endpoint, "elasticsearch", clusterID.String())

			// println(tskCtl.Endpoint)
			tskCtl.Version = latestVersion(versions)
			tskCtl.StorePath = filepath.Join(tskCtl.StorePath, "elasticsearch", clusterID.String(), "elasticsearch")

			for _, task := range tasks {
				// Add task to wait group
				wg.Add(1)

				// task.Filename = filepath.Join("elasticsearch/{{ .cluster_id }}", task.Filename)
				// task.Filename = filepath.Join("elasticsearch/{{ .cluster_id }}", task.Filename)

				task.Headers = eceHeader()
				tskCtl.Task = task

				go tskCtl.TaskExecuteWithWaitGroup(&wg)
			}

			wg.Wait()
			clusterWait.Done()

		}(*elasticsearch.NewElasticsearchAPIset(), t, &clusterWait)

		return true // keep iterating
	})

	// wait for each cluster
	clusterWait.Wait()
}

func latestVersion(versions gjson.Result) string {
	latest, _ := semver.NewVersion("0")
	for _, ver := range versions.Array() {
		// setup check for greater than the current latest version value
		c, _ := semver.NewConstraint("> " + latest.String())
		v, _ := semver.NewVersion(ver.String())

		// greater than returns true, update the latest version
		if c.Check(v) {
			latest = v
		}
	}
	return latest.String()
	// TODO: if anything goes wrong in the version checking, this should return nil?
	// the nil value could be used to just apply a generic set of API commands.
}

func eceHeader() map[string]string {
	header := map[string]string{"X-Management-Request": "true"}
	return header
}
