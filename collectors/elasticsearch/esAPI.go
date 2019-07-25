package elasticsearch

import (
	"net/http"

	"github.com/elastic/ece-support-diagnostics/helpers"
)

// ECE should prepend `elasticsearch/{{ .ClusterID }}/` to the URI

func NewElasticsearchAPIset() *helpers.Tasks {
	APItasks := helpers.Tasks{
		// Common
		helpers.Task{
			Filename: "/alias.json",
			Uri:      "/_alias?pretty&human",
			Method:   http.MethodGet,
		},
		helpers.Task{
			Filename: "/cat_aliases.txt",
			Uri:      "/_cat/aliases?v",
		},
		helpers.Task{
			Filename: "/cat_allocation.txt",
			Uri:      "/_cat/allocation?v",
		},
		helpers.Task{
			Filename: "/cat_count.txt",
			Uri:      "/_cat/count",
		},
		helpers.Task{
			Filename: "/cat_fielddata.txt",
			Uri:      "/_cat/fielddata?v",
		},
		helpers.Task{
			Filename: "/cat_health.txt",
			Uri:      "/_cat/health?v"},
		helpers.Task{
			Filename: "/cat_indices.txt",
			Uri:      "/_cat/indices?v",
			Versions: "< 5.2.x",
		},
		helpers.Task{
			Filename: "/cat_master.txt",
			Uri:      "/_cat/master?format=json",
		},
		helpers.Task{
			Filename: "/cat_nodes.txt",
			Uri:      "/_cat/nodes?v&h=n,m,i,r,d,hp,rp,cpu,load_1m,load_5m,load_15m,nodeId",
		},
		helpers.Task{
			Filename: "/cat_pending_tasks.txt",
			Uri:      "/_cat/pending_tasks?v",
		},
		helpers.Task{
			Filename: "/cat_segments.txt",
			Uri:      "/_cat/segments?v",
		},
		helpers.Task{
			Filename: "/cat_recovery.txt",
			Uri:      "/_cat/recovery?v",
		},
		helpers.Task{
			Filename: "/cat_shards.txt",
			Uri:      "/_cat/shards?v",
		},
		helpers.Task{
			Filename: "/cat_thread_pool.txt",
			Uri:      "/_cat/thread_pool?v",
		},
		helpers.Task{
			Filename: "/cluster_health.json",
			Uri:      "/_cluster/health?pretty",
		},
		helpers.Task{
			Filename: "/cluster_pending_tasks.json",
			Uri:      "/_cluster/pending_tasks?pretty&human",
		},
		helpers.Task{
			Filename: "/cluster_settings.json",
			Uri:      "/_cluster/settings?pretty&flat_settings",
		},
		helpers.Task{
			Filename: "/cluster_state.json",
			Uri:      "/_cluster/state?pretty&human",
		},
		helpers.Task{
			Filename: "/cluster_stats.json",
			Uri:      "/_cluster/stats?pretty&human",
		},
		helpers.Task{
			Filename: "/fielddata.json",
			Uri:      "/_cat/fielddata?format=json&bytes&pretty",
		},
		helpers.Task{
			Filename: "/fielddata_stats.json",
			Uri:      "/_nodes/stats/indices/fielddata?pretty=true&fields=*",
		},
		helpers.Task{
			Filename: "/indices_stats.json",
			Uri:      "/_stats?level=shards&pretty&human",
		},
		helpers.Task{
			Filename: "/mapping.json",
			Uri:      "/_mapping?pretty",
		},
		helpers.Task{
			Filename: "/master.json",
			Uri:      "/_cat/master?format=json",
		},
		helpers.Task{
			Filename: "/nodes_hot_threads.txt",
			Uri:      "/_nodes/hot_threads?threads=10000",
		},
		helpers.Task{
			Filename: "/nodes_stats.json",
			Uri:      "/_nodes/stats?pretty&human",
		},
		helpers.Task{
			Filename: "/nodes.json",
			Uri:      "/_nodes?pretty&human",
		},
		helpers.Task{
			Filename: "/plugins.json",
			Uri:      "/_cat/plugins?format=json",
		},
		helpers.Task{
			Filename: "/recovery.json",
			Uri:      "/_recovery?pretty&human&detailed=true",
		},
		helpers.Task{
			Filename: "/segments.json",
			Uri:      "/_segments?pretty&human",
		},
		helpers.Task{
			Filename: "/settings.json",
			Uri:      "/_settings?pretty&human",
		},
		helpers.Task{
			Filename: "/shards.json",
			Uri:      "/_cat/shards?format=json&bytes=b&pretty",
			Callback: unassignedShard{},
		},
		helpers.Task{
			Filename: "/templates.json",
			Uri:      "/_template?pretty",
		},
		helpers.Task{
			Filename: "/version.json",
			Uri:      "/",
		},
		helpers.Task{
			Filename: "/watcher_stats.json",
			Uri:      "/_watcher/stats/_all",
		},
		helpers.Task{
			Filename: "/watcher_stack.json",
			Uri:      "/_watcher/stats?emit_stacktraces=true",
		},

		// 1.x
		helpers.Task{
			Filename: "/licenses.json",
			Uri:      "/_licenses",
			Versions: "^1.x",
		},
		// 2.x
		helpers.Task{
			Filename: "/cat_nodeattrs.txt",
			Uri:      "/_cat/nodeattrs?v&h=node,id,pid,host,ip,port,attr,value",
			Versions: ">= 2.x",
		},
		helpers.Task{
			Filename: "/cat_repositories.txt",
			Uri:      "/_cat/repositories?v",
			Versions: ">= 2.x",
		},
		helpers.Task{
			Filename: "/count.json",
			Uri:      "/_count",
			Versions: ">= 2.x",
		},
		helpers.Task{
			Filename: "/licenses.json",
			Uri:      "/_license?pretty",
			Versions: ">= 2.x",
		},
		helpers.Task{
			Filename: "/secUrity_users.json",
			Uri:      "/_shield/user?pretty",
			Versions: ">= 2.x",
		},
		helpers.Task{
			Filename: "/secUrity_roles.json",
			Uri:      "/_shield/role?pretty",
			Versions: ">= 2.x",
		},
		helpers.Task{
			Filename: "/shard_stores.json",
			Uri:      "/_shard_stores?pretty",
			Versions: ">= 2.x",
		},
		helpers.Task{
			Filename: "/tasks.json",
			Uri:      "/_tasks?pretty&human&detailed=true",
			Versions: ">= 2.x",
		},
		// 5.0 ?
		helpers.Task{
			Filename: "/allocation_explain.json",
			Uri:      "/_cluster/allocation/explain?pretty",
			Versions: ">= 5.x",
		},
		helpers.Task{
			Filename: "/allocation_explain_disk.json",
			Uri:      "/_cluster/allocation/explain?include_disk_info=true&pretty",
			Versions: ">= 5.x",
		},
		helpers.Task{
			Filename: "/fielddata_stats.json",
			Uri:      "/_nodes/stats/indices/fielddata?level=shards&pretty=true&fields=*",
			Versions: ">= 5.x",
		},
		helpers.Task{
			Filename: "/ml_anomaly_detectors.json",
			Uri:      "/_xpack/ml/anomaly_detectors?pretty",
			Versions: ">= 5.x",
		},
		helpers.Task{
			Filename: "/ml_datafeeds.json",
			Uri:      "/_xpack/ml/datafeeds?pretty",
			Versions: ">= 5.x",
		},
		helpers.Task{
			Filename: "/ml_stats.json",
			Uri:      "/_xpack/ml/anomaly_detectors/_stats?pretty",
			Versions: ">= 5.x",
		},
		helpers.Task{
			Filename: "/pipelines.json",
			Uri:      "/_ingest/pipeline/*?pretty&human",
			Versions: ">= 5.x",
		},
		helpers.Task{
			Filename: "/secUrity_users.json",
			Uri:      "/_xpack/secUrity/user?pretty",
			Versions: ">= 5.x",
		},
		helpers.Task{
			Filename: "/secUrity_roles.json",
			Uri:      "/_xpack/secUrity/role?pretty",
			Versions: ">= 5.x",
		},
		helpers.Task{
			Filename: "/secUrity_role_mappings.json",
			Uri:      "/_xpack/secUrity/role_mapping?pretty",
			Versions: ">= 5.x",
		},
		helpers.Task{
			Filename: "/xpack.json",
			Uri:      "/_xpack/usage?pretty&human",
			Versions: ">= 5.x",
		},
		// 5.2 ?
		helpers.Task{
			Filename: "/cat_indices.txt",
			Uri:      "/_cat/indices?v&s=index",
			Versions: ">= 5.2.x",
		},
		helpers.Task{
			Filename: "/cat_segments.txt",
			Uri:      "/_cat/segments?v&s=index",
			Versions: ">= 5.2.x",
		},
		helpers.Task{
			Filename: "/cat_shards.txt",
			Uri:      "/_cat/shards?v&s=index",
			Versions: ">= 5.2.x",
		},
		// 6.0 ?
		helpers.Task{
			Filename: "/nodes_usage.json",
			Uri:      "/_nodes/usage?pretty",
			Versions: ">= 6.x",
		},
		helpers.Task{
			Filename: "/remote_cluster_info.json",
			Uri:      "/_remote/info",
			Versions: ">= 6.x",
		},
		// 6.3 ?
		helpers.Task{
			Filename: "/rollup_jobs.json",
			Uri:      "/_xpack/rollup/job/_all",
			Versions: ">= 6.3.x",
		},
		helpers.Task{
			Filename: "/rollup_caps.json",
			Uri:      "/_xpack/rollup/data/_all",
			Versions: ">= 6.3.x",
		},
		// 6.4 ?
		helpers.Task{
			Filename: "/cluster_settings_defaults.json",
			Uri:      "/_cluster/settings?include_defaults&pretty&flat_settings",
			Versions: ">= 6.4.x",
		},
		// 6.5?
		helpers.Task{
			Filename: "/rollup_index_caps.json",
			Uri:      "/*/_xpack/rollup/data",
			Versions: ">= 6.5.x",
		},
		helpers.Task{
			Filename: "/secUrity_priv.json",
			Uri:      "/_xpack/secUrity/privilege?pretty",
			Versions: ">= 6.5.x",
		},
		helpers.Task{
			Filename: "/ccr_stats.json",
			Uri:      "/_ccr/stats?pretty",
			Versions: ">= 6.5.x",
		},
		helpers.Task{
			Filename: "/ccr_autofollow_patterns.json",
			Uri:      "/_ccr/auto_follow?pretty",
			Versions: ">= 6.5.x",
		},
		// 6.6 ?
		helpers.Task{
			Filename: "/ilm_explain.json",
			Uri:      "/*/_ilm/explain?human&pretty",
			Versions: ">= 6.6.x",
		},
		helpers.Task{
			Filename: "/ilm_policies.json",
			Uri:      "/_ilm/policy?human&pretty",
			Versions: ">= 6.6.x",
		},
		helpers.Task{
			Filename: "/ilm_status.json",
			Uri:      "/_ilm/status?pretty",
			Versions: ">= 6.6.x",
		},
		// 6.7 ?
		helpers.Task{
			Filename: "/ccr_follower_info.json",
			Uri:      "/_all/_ccr/info?pretty",
			Versions: ">= 6.6.x",
		},
	}
	return &APItasks

}
