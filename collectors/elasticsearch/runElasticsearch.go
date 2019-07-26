package elasticsearch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/elastic/ece-support-diagnostics/config"
	"github.com/elastic/ece-support-diagnostics/helpers"
	"github.com/tidwall/gjson"
)

func exampleElasticsearchRunner(cfg *config.Config) {

	// store := new(fileStore)

	// Used for templating and version checking.
	es := helpers.TaskContext{
		Version:  "5.6",
		Endpoint: "http://localhost:9200/",
		User:     "elastic",
		Pass:     "changeme",
		Client:   cfg.HTTPclient,
		Store:    cfg.Store,
	}

	// loop over each Task and dispatch
	func(tasks helpers.Tasks, tskCtl helpers.TaskContext) {

		var wg sync.WaitGroup

		for _, task := range tasks {
			// Add task to wait group
			wg.Add(1)

			es.Task = task

			go es.TaskExecuteWithWaitGroup(&wg)
		}

		wg.Wait()

	}(*NewElasticsearchAPIset(), es)
}

type unassignedShard struct {
	Index   string `json:"index"`
	Shard   int64  `json:"shard"`
	Primary bool   `json:"primary"`
}

func (u unassignedShard) Exec(t helpers.TaskContext, p []byte) {
	checkUnassignedShards(t, p)
}

func checkUnassignedShards(t helpers.TaskContext, payload []byte) {
	var wg sync.WaitGroup

	// Get all the unassigned shards and loop over them
	unassignedShards := gjson.GetBytes(payload, `#(state!="STARTED")#`)
	unassignedShards.ForEach(func(key, value gjson.Result) bool {
		// println(value.String())
		index := gjson.Get(value.Raw, "index")
		shard := gjson.Get(value.Raw, "shard")
		prirep := gjson.Get(value.Raw, "prirep")

		primary := false
		if prirep.String() == "p" {
			primary = true
		}
		Shard := unassignedShard{
			Index:   index.String(),
			Shard:   shard.Int(),
			Primary: primary,
		}

		b, err := json.Marshal(Shard)
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Println(string(b))

		// Explain
		wg.Add(1)
		t.Task = helpers.Task{
			Filename:    fmt.Sprintf("unallocated_shards/%s_%s%s_allocation_explain.json", index, shard, prirep),
			Uri:         "/_cluster/allocation/explain?pretty",
			Versions:    ">= 5.x",
			Method:      http.MethodPost,
			RequestBody: b,
			Headers:     postHeaderJSON(),
		}
		go t.TaskExecuteWithWaitGroup(&wg)

		// Disk Info
		wg.Add(1)
		t.Task = helpers.Task{
			Filename:    fmt.Sprintf("unallocated_shards/%s_%s%s_allocation_explain_disk.json", index, shard, prirep),
			Uri:         "/_cluster/allocation/explain?include_disk_info=true&pretty",
			Versions:    ">= 5.x",
			Method:      http.MethodPost,
			RequestBody: b,
			Headers:     postHeaderJSON(),
		}
		go t.TaskExecuteWithWaitGroup(&wg)

		return true // keep iterating
	})

	wg.Wait()
}

func postHeaderJSON() map[string]string {
	header := map[string]string{"Content-Type": "application/json", "X-Management-Request": "true"}
	return header
}
