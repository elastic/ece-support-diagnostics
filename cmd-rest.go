package ecediag

import (
	"strings"

	"github.com/elastic/ece-support-diagnostics/collectors/restAPI"
)

// this is a helper function to avoid long repetitive paths for a large number of Rest calls
func eceProxy(uri string) string {
	eceBase := "/api/v1/clusters/elasticsearch/{{ .cluster_id }}/proxy/"
	return eceBase + strings.TrimLeft(uri, "/")
}

func eceHeader() map[string]string {
	header := map[string]string{"X-Management-Request": "true"}
	return header
}

// If `Sub` is present, the JSON response will be unpacked to a map and provided
//  to the sub items for templating. `Loop` will specify the path to an array
//  to iterate on. `Loop` will be looped for the `Sub` items.
var rest = []restAPI.Rest{
	restAPI.Rest{
		Filename: "ece/platform.json",
		URI:      "/api/v1/platform",
		Headers:  eceHeader(),
	},

	restAPI.Rest{
		Filename: "ece/allocators.json",
		URI:      "/api/v1/platform/infrastructure/allocators",
		Headers:  eceHeader(),
	},
	restAPI.Rest{
		Filename: "ece/runners.json",
		URI:      "/api/v1/platform/infrastructure/runners",
		Headers:  eceHeader(),
	},
	restAPI.Rest{
		Filename: "ece/proxies.json",
		URI:      "/api/v1/platform/infrastructure/proxies",
		Headers:  eceHeader(),
	},
	// TODO: Add kibana api/stats, api/settings, api/status ?
	restAPI.Rest{
		Filename: "ece/kibana_clusters.json",
		URI:      "/api/v1/clusters/kibana",
		Headers:  eceHeader(),
	},
	restAPI.Rest{
		Filename: "ece/es_clusters.json",
		URI:      "/api/v1/clusters/elasticsearch",
		Loop:     "elasticsearch_clusters",
		Headers:  eceHeader(),
		Sub: []restAPI.Rest{
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/_ece_cluster_info.json",
				URI:      "/api/v1/clusters/elasticsearch/{{ .cluster_id }}",
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/_ece_plan.json",
				URI:      "/api/v1/clusters/elasticsearch/{{ .cluster_id }}/plan/activity",
				Headers:  eceHeader(),
			},
			// restAPI.Rest{
			// 	// THIS MUST BE COLLECTED AS A NON-PRIVILEGED USER!!!
			// 	// if not, it would expose sensitive data that we should not collect.
			// 	Filename: "elasticsearch/{{ .cluster_id }}/_ece_cluster_metadata.json",
			// 	URI:  "/api/v1/clusters/elasticsearch/{{ .cluster_id }}/metadata/raw",
			// },
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/alias.json",
				URI:      eceProxy("/_alias?pretty&human"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_aliases.txt",
				URI:      eceProxy("_cat/aliases?v"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_tasks.txt",
				URI:      eceProxy("_cat/tasks"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_allocation.txt",
				URI:      eceProxy("_cat/allocation?v"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_count.txt",
				URI:      eceProxy("_cat/count"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_fielddata.txt",
				URI:      eceProxy("_cat/fielddata?v"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_health.txt",
				URI:      eceProxy("_cat/health?v"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_indices.txt",
				URI:      eceProxy("_cat/indices?v&s=name"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_master.txt",
				URI:      eceProxy("_cat/master?format=json"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_nodes.txt",
				URI:      eceProxy("_cat/nodes?v&h=n,m,i,r,d,hp,rp,cpu,load_1m,load_5m,load_15m,nodeId"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_pending_tasks.txt",
				URI:      eceProxy("_cat/pending_tasks?v"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_segments.txt",
				URI:      eceProxy("_cat/segments?v"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_recovery.txt",
				URI:      eceProxy("_cat/recovery?v"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_shards.txt",
				URI:      eceProxy("_cat/shards?v"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_thread_pool.txt",
				URI:      eceProxy("_cat/thread_pool?v"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cluster_health.json",
				URI:      eceProxy("_cluster/health?pretty"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cluster_pending_tasks.json",
				URI:      eceProxy("_cluster/pending_tasks?pretty&human"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cluster_settings.json",
				URI:      eceProxy("_cluster/settings?pretty&flat_settings"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cluster_state.json",
				URI:      eceProxy("_cluster/state?pretty&human"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cluster_stats.json",
				URI:      eceProxy("_cluster/stats?pretty&human"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/fielddata.txt",
				URI:      eceProxy("_cat/fielddata?format=json&bytes&pretty"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/fielddata_stats.json",
				URI:      eceProxy("_nodes/stats/indices/fielddata?pretty=true&fields=*"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/indices_stats.json",
				URI:      eceProxy("_stats?pretty&human"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/mapping.json",
				URI:      eceProxy("_mapping?pretty"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/master.json",
				URI:      eceProxy("_cat/master?format=json"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/nodes_hot_threads.txt",
				URI:      eceProxy("_nodes/hot_threads?threads=10000"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/nodes_stats.json",
				URI:      eceProxy("_nodes/stats?pretty&human"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/nodes.json",
				URI:      eceProxy("_nodes?pretty&human"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/plugins.json",
				URI:      eceProxy("_cat/plugins?format=json"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/recovery.json",
				URI:      eceProxy("_recovery?pretty&human&detailed=true"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/segments.json",
				URI:      eceProxy("_segments?pretty&human&verbose"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/settings.json",
				URI:      eceProxy("_settings?pretty&human"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/shards.json",
				URI:      eceProxy("_cat/shards?format=json&bytes=b&pretty"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/templates.json",
				URI:      eceProxy("_template?pretty"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/version.json",
				URI:      eceProxy(""),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/watcher_stats.json",
				URI:      eceProxy("_watcher/stats/_all"),
				Headers:  eceHeader(),
			},
			restAPI.Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/watcher_stack.json",
				URI:      eceProxy("_watcher/stats?emit_stacktraces=true"),
				Headers:  eceHeader(),
			},
		},
	},
}
