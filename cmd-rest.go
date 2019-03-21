package ecediag

import "strings"

// Rest defines HTTP requests to be collected
type Rest struct {
	Method   string
	Request  string
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
		Request:  "/api/v1/platform",
		Filename: "ece/platform.json",
	},
	Rest{
		Request:  "/api/v1/platform/infrastructure/allocators",
		Filename: "ece/allocators.json",
	},
	Rest{
		Request:  "/api/v1/platform/infrastructure/runners",
		Filename: "ece/runners.json",
	},
	Rest{
		Request:  "/api/v1/platform/infrastructure/proxies",
		Filename: "ece/proxies.json",
	},
	Rest{
		Request:  "/api/v1/clusters/elasticsearch",
		Filename: "ece/es_clusters.json",
		Loop:     "elasticsearch_clusters",
		Sub: []Rest{
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/_ece_cluster_info.json",
				Request:  "/api/v1/clusters/elasticsearch/{{ .cluster_id }}",
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/_ece_plan.json",
				Request:  "/api/v1/clusters/elasticsearch/{{ .cluster_id }}/plan/activity",
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/alias.json",
				Request:  eceProxy("/_alias?pretty&human"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_aliases.txt",
				Request:  eceProxy("_cat/aliases?v"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_tasks.txt",
				Request:  eceProxy("_cat/tasks"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_allocation.txt",
				Request:  eceProxy("_cat/allocation?v"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_count.txt",
				Request:  eceProxy("_cat/count"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_fielddata.txt",
				Request:  eceProxy("_cat/fielddata?v"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_health.txt",
				Request:  eceProxy("_cat/health?v"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_indices.txt",
				Request:  eceProxy("_cat/indices?v&s=name"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_master.txt",
				Request:  eceProxy("_cat/master?format=json"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_nodes.txt",
				Request:  eceProxy("_cat/nodes?v&h=n,m,i,r,d,hp,rp,cpu,load_1m,load_5m,load_15m,nodeId"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_pending_tasks.txt",
				Request:  eceProxy("_cat/pending_tasks?v"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_segments.txt",
				Request:  eceProxy("_cat/segments?v"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_recovery.txt",
				Request:  eceProxy("_cat/recovery?v"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_shards.txt",
				Request:  eceProxy("_cat/shards?v"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cat_thread_pool.txt",
				Request:  eceProxy("_cat/thread_pool?v"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cluster_health.json",
				Request:  eceProxy("_cluster/health?pretty"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cluster_pending_tasks.json",
				Request:  eceProxy("_cluster/pending_tasks?pretty&human"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cluster_settings.json",
				Request:  eceProxy("_cluster/settings?pretty&flat_settings"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cluster_state.json",
				Request:  eceProxy("_cluster/state?pretty&human"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/cluster_stats.json",
				Request:  eceProxy("_cluster/stats?pretty&human"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/fielddata.txt",
				Request:  eceProxy("_cat/fielddata?format=json&bytes&pretty"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/fielddata_stats.json",
				Request:  eceProxy("_nodes/stats/indices/fielddata?pretty=true&fields=*"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/indices_stats.json",
				Request:  eceProxy("_stats?pretty&human"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/mapping.json",
				Request:  eceProxy("_mapping?pretty"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/master.json",
				Request:  eceProxy("_cat/master?format=json"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/nodes_hot_threads.txt",
				Request:  eceProxy("_nodes/hot_threads?threads=10000"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/nodes_stats.json",
				Request:  eceProxy("_nodes/stats?pretty&human"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/nodes.json",
				Request:  eceProxy("_nodes?pretty&human"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/plugins.json",
				Request:  eceProxy("_cat/plugins?format=json"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/recovery.json",
				Request:  eceProxy("_recovery?pretty&human&detailed=true"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/segments.json",
				Request:  eceProxy("_segments?pretty&human&verbose"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/settings.json",
				Request:  eceProxy("_settings?pretty&human"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/shards.json",
				Request:  eceProxy("_cat/shards?format=json&bytes=b&pretty"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/templates.json",
				Request:  eceProxy("_template?pretty"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/version.json",
				Request:  eceProxy(""),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/watcher_stats.json",
				Request:  eceProxy("_watcher/stats/_all"),
			},
			Rest{
				Filename: "elasticsearch/{{ .cluster_id }}/watcher_stack.json",
				Request:  eceProxy("_watcher/stats?emit_stacktraces=true"),
			},
		},
	},
}
