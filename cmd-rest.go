package ecediag

import "strings"

// Rest defines HTTP URIs to be collected
type Rest struct {
	Method   string
	URI      string
	Filename string
	Loop     string
	Sub      []Rest
}

// this is a helper function to avoid long repetitive paths for a large number of Rest calls
func eceProxy(uri string) string {
	eceBase := "/api/v1/clusters/elasticsearch/{{ .cluster_id }}/proxy/"
	return eceBase + strings.TrimLeft(uri, "/")
}

// If `Sub` is present, the JSON response will be unpacked to a map and provided
//  to the sub items for templating. `Loop` will specify the path to an array
//  to iterate on. `Loop` will be looped for the `Sub` items.
var rest = []Rest{
	Rest{
		Filename: "ece/platform.json",
		URI:      "/api/v1/platform",
	},
	Rest{
		Filename: "ece/allocators.json",
		URI:      "/api/v1/platform/infrastructure/allocators",
	},
	Rest{
		Filename: "ece/runners.json",
		URI:      "/api/v1/platform/infrastructure/runners",
	},
	Rest{
		Filename: "ece/proxies.json",
		URI:      "/api/v1/platform/infrastructure/proxies",
	},
	Rest{
		Filename: "ece/es_clusters.json",
		URI:      "/api/v1/clusters/elasticsearch",
		Loop:     "elasticsearch_clusters",
		Sub: []Rest{
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/_ece_cluster_info.json",
				URI:      "/api/v1/clusters/elasticsearch/{{ .cluster_id }}",
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/_ece_plan.json",
				URI:      "/api/v1/clusters/elasticsearch/{{ .cluster_id }}/plan/activity",
			},
			// Rest{
			// 	// THIS MUST BE COLLECTED AS A NON-PRIVILEGED USER!!!
			// 	// if not, it would expose sensitive data that we should not collect.
			// 	Filename: "elasticsearch/{{ .cluster_id }}/_ece_cluster_metadata.json",
			// 	URI:  "/api/v1/clusters/elasticsearch/{{ .cluster_id }}/metadata/raw",
			// },
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/alias.json",
				URI:      eceProxy("/_alias?pretty&human"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_aliases.txt",
				URI:      eceProxy("_cat/aliases?v"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_tasks.txt",
				URI:      eceProxy("_cat/tasks"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_allocation.txt",
				URI:      eceProxy("_cat/allocation?v"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_count.txt",
				URI:      eceProxy("_cat/count"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_fielddata.txt",
				URI:      eceProxy("_cat/fielddata?v"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_health.txt",
				URI:      eceProxy("_cat/health?v"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_indices.txt",
				URI:      eceProxy("_cat/indices?v&s=name"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_master.txt",
				URI:      eceProxy("_cat/master?format=json"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_nodes.txt",
				URI:      eceProxy("_cat/nodes?v&h=n,m,i,r,d,hp,rp,cpu,load_1m,load_5m,load_15m,nodeId"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_pending_tasks.txt",
				URI:      eceProxy("_cat/pending_tasks?v"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_segments.txt",
				URI:      eceProxy("_cat/segments?v"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_recovery.txt",
				URI:      eceProxy("_cat/recovery?v"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_shards.txt",
				URI:      eceProxy("_cat/shards?v"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_thread_pool.txt",
				URI:      eceProxy("_cat/thread_pool?v"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cluster_health.json",
				URI:      eceProxy("_cluster/health?pretty"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cluster_pending_tasks.json",
				URI:      eceProxy("_cluster/pending_tasks?pretty&human"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cluster_settings.json",
				URI:      eceProxy("_cluster/settings?pretty&flat_settings"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cluster_state.json",
				URI:      eceProxy("_cluster/state?pretty&human"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cluster_stats.json",
				URI:      eceProxy("_cluster/stats?pretty&human"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/fielddata.txt",
				URI:      eceProxy("_cat/fielddata?format=json&bytes&pretty"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/fielddata_stats.json",
				URI:      eceProxy("_nodes/stats/indices/fielddata?pretty=true&fields=*"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/indices_stats.json",
				URI:      eceProxy("_stats?pretty&human"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/mapping.json",
				URI:      eceProxy("_mapping?pretty"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/master.json",
				URI:      eceProxy("_cat/master?format=json"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/nodes_hot_threads.txt",
				URI:      eceProxy("_nodes/hot_threads?threads=10000"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/nodes_stats.json",
				URI:      eceProxy("_nodes/stats?pretty&human"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/nodes.json",
				URI:      eceProxy("_nodes?pretty&human"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/plugins.json",
				URI:      eceProxy("_cat/plugins?format=json"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/recovery.json",
				URI:      eceProxy("_recovery?pretty&human&detailed=true"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/segments.json",
				URI:      eceProxy("_segments?pretty&human&verbose"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/settings.json",
				URI:      eceProxy("_settings?pretty&human"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/shards.json",
				URI:      eceProxy("_cat/shards?format=json&bytes=b&pretty"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/templates.json",
				URI:      eceProxy("_template?pretty"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/version.json",
				URI:      eceProxy(""),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/watcher_stats.json",
				URI:      eceProxy("_watcher/stats/_all"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/watcher_stack.json",
				URI:      eceProxy("_watcher/stats?emit_stacktraces=true"),
			},
		},
	},
}
