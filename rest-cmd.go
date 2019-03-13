package ecediag

import "strings"

type Rest struct {
	Method    string
	Request   string
	Filename  string
	WithItems string
	Sub       []Rest
}

func eceProxy(uri string) string {
	eceBase := "/api/v1/clusters/elasticsearch/{{ .cluster_id }}/proxy/"
	return eceBase + strings.TrimLeft(uri, "/")
}

// If `Collector` is present, the JSON response will be unmarshalled to it.
//  Any request in `Sub` will be provided the Collect struct for text substituion
var rest = []Rest{
	Rest{
		Request:  "/api/v1/platform/infrastructure/allocators",
		Filename: "allocators.json",
	},
	Rest{
		Request:  "/api/v1/platform/infrastructure/runners",
		Filename: "runners.json",
	},
	Rest{
		Request:  "/api/v1/platform/infrastructure/proxies",
		Filename: "proxies.json",
	},
	Rest{
		Filename:  "es_clusters.json",
		Request:   "/api/v1/clusters/elasticsearch",
		WithItems: "elasticsearch_clusters",
		Sub: []Rest{
			Rest{
				Filename: "{{ .cluster_id }}/ece/plan.json",
				Request:  "/api/v1/clusters/elasticsearch/{{ .cluster_id }}/plan/activity",
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/alias.json",
				Request:  eceProxy("/_alias?pretty&human"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/cat_aliases.txt",
				Request:  eceProxy("_cat/aliases?v"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/cat_allocation.txt",
				Request:  eceProxy("_cat/allocation?v"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/cat_count.txt",
				Request:  eceProxy("_cat/count"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/cat_fielddata.txt",
				Request:  eceProxy("_cat/fielddata?v"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/cat_health.txt",
				Request:  eceProxy("_cat/health?v"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/cat_indices.txt",
				Request:  eceProxy("_cat/indices?v"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/cat_master.txt",
				Request:  eceProxy("_cat/master?format=json"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/cat_nodes.txt",
				Request:  eceProxy("_cat/nodes?v&h=n,m,i,r,d,hp,rp,cpu,load_1m,load_5m,load_15m,nodeId"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/cat_pending_tasks.txt",
				Request:  eceProxy("_cat/pending_tasks?v"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/cat_segments.txt",
				Request:  eceProxy("_cat/segments?v"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/cat_recovery.txt",
				Request:  eceProxy("_cat/recovery?v"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/cat_shards.txt",
				Request:  eceProxy("_cat/shards?v"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/cat_thread_pool.txt",
				Request:  eceProxy("_cat/thread_pool?v"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/cluster_health.json",
				Request:  eceProxy("_cluster/health?pretty"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/cluster_pending_tasks.json",
				Request:  eceProxy("_cluster/pending_tasks?pretty&human"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/cluster_settings.json",
				Request:  eceProxy("_cluster/settings?pretty&flat_settings"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/cluster_state.json",
				Request:  eceProxy("_cluster/state?pretty&human"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/cluster_stats.json",
				Request:  eceProxy("_cluster/stats?pretty&human"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/fielddata.txt",
				Request:  eceProxy("_cat/fielddata?format=json&bytes&pretty"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/fielddata_stats.json",
				Request:  eceProxy("_nodes/stats/indices/fielddata?pretty=true&fields=*"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/indices_stats.json",
				Request:  eceProxy("_stats?pretty&human"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/mapping.json",
				Request:  eceProxy("_mapping?pretty"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/master.json",
				Request:  eceProxy("_cat/master?format=json"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/nodes_hot_threads.txt",
				Request:  eceProxy("_nodes/hot_threads?threads=10000"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/nodes_stats.json",
				Request:  eceProxy("_nodes/stats?pretty&human"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/nodes.json",
				Request:  eceProxy("_nodes?pretty&human"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/plugins.json",
				Request:  eceProxy("_cat/plugins?format=json"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/recovery.json",
				Request:  eceProxy("_recovery?pretty&human&detailed=true"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/segments.json",
				Request:  eceProxy("_segments?pretty&human&verbose"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/settings.json",
				Request:  eceProxy("_settings?pretty&human"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/shards.json",
				Request:  eceProxy("_cat/shards?format=json&bytes=b&pretty"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/templates.json",
				Request:  eceProxy("_template?pretty"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/version.json",
				Request:  eceProxy(""),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/watcher_stats.json",
				Request:  eceProxy("_watcher/stats/_all"),
			},
			Rest{
				Filename: "{{ .cluster_id }}/es/watcher_stack.json",
				Request:  eceProxy("_watcher/stats?emit_stacktraces=true"),
			},
		},
	},
}
